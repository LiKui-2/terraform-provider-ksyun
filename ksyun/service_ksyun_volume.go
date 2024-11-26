package ksyun

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

type EbsService struct {
	client *KsyunClient
}

func (s *EbsService) ReadVolumes(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	return pageQuery(condition, "MaxResults", "Marker", 50, 0, func(condition map[string]interface{}) ([]interface{}, error) {
		conn := s.client.ebsconn
		action := "DescribeVolumes"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = conn.DescribeVolumes(nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = conn.DescribeVolumes(&condition)
			if err != nil {
				return data, err
			}
		}

		results, err = getSdkValue("Volumes", *resp)
		if err != nil {
			return data, err
		}
		data = results.([]interface{})
		return data, err
	})
}

func (s *EbsService) ReadVolume(d *schema.ResourceData, volumeId string, allProject bool) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if volumeId == "" {
		volumeId = d.Id()
	}
	req := map[string]interface{}{
		"VolumeId.1": volumeId,
	}
	if allProject {
		err = addProjectInfoAll(d, &req, s.client)
		if err != nil {
			return data, err
		}
	} else {
		err = addProjectInfo(d, &req, s.client)
		if err != nil {
			return data, err
		}
	}

	results, err = s.ReadVolumes(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Volume %s not exist ", volumeId)
	}
	return data, err
}

func (s *EbsService) volumeStateRefreshFunc(d *schema.ResourceData, volumeId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			err error
		)
		data, err := s.ReadVolume(d, volumeId, true)
		if err != nil {
			return nil, "", err
		}

		status, err := getSdkValue("VolumeStatus", data)
		if err != nil {
			return nil, "", err
		}

		for _, v := range failStates {
			if v == status.(string) {
				return nil, "", fmt.Errorf("volume status  error, status:%v", status)
			}
		}
		return data, status.(string), nil
	}
}

func (s *EbsService) checkVolumeState(d *schema.ResourceData, volumeId string, target []string, timeout time.Duration) (state interface{}, err error) {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     target,
		Refresh:    s.volumeStateRefreshFunc(d, volumeId, []string{"error"}),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 1 * time.Minute,
	}
	return stateConf.WaitForState()
}

func (s *EbsService) ReadAndSetVolume(d *schema.ResourceData, r *schema.Resource) (err error) {
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		data, callErr := s.ReadVolume(d, "", false)
		if callErr != nil {
			if !d.IsNewResource() {
				return resource.NonRetryableError(callErr)
			}
			if notFoundError(callErr) {
				return resource.RetryableError(callErr)
			} else {
				return resource.NonRetryableError(fmt.Errorf("error on  reading volume %q, %s", d.Id(), callErr))
			}
		} else {
			extra := chargeExtraForVpc(data)
			extra["Tags"] = SdkResponseMapping{
				Field: "tags",
				FieldRespFunc: func(i interface{}) interface{} {
					tags := i.([]interface{})
					tagMap := make(map[string]interface{})
					for _, tag := range tags {
						_m := tag.(map[string]interface{})
						tagMap[_m["TagKey"].(string)] = _m["TagValue"].(string)
					}
					return tagMap
				},
			}

			SdkResponseAutoResourceData(d, r, data, extra)
			return nil
		}
	})
}

func (s *EbsService) ReadAndSetVolumes(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "VolumeId",
			Type:    TransformWithN,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.ReadVolumes(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		idFiled:     "VolumeId",
		targetField: "volumes",
		extra:       map[string]SdkResponseMapping{},
	})
}

func (s *EbsService) CreateVolumeCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"online_resize": {Ignore: true},
		"tags":          {Ignore: true},
	}
	req, err := SdkRequestAutoMapping(d, r, false, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})
	if err != nil {
		return callback, err
	}
	callback = ApiCall{
		param:  &req,
		action: "CreateVolume",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.ebsconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateVolume(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("VolumeId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			_, err = s.checkVolumeState(d, "", []string{"available"}, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return err
			}
			return d.Set("online_resize", d.Get("online_resize"))
		},
	}
	return callback, err
}

func (s *EbsService) CreateVolume(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.CreateVolumeCall(d, r)
	if err != nil {
		return err
	}
	tagsService := TagService{client: s.client}
	tagsCall, err := tagsService.ReplaceResourcesTagsWithResourceCall(d, r, "volume", false, false)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call, tagsCall}, d, s.client, true)
}

func (s *EbsService) ModifyVolumeProjectCall(d *schema.ResourceData, resource *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"project_id": {},
	}
	updateReq, err := SdkRequestAutoMapping(d, resource, true, transform, nil)
	if err != nil {
		return callback, err
	}
	if len(updateReq) > 0 {
		callback = ApiCall{
			param: &updateReq,
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				return resp, ModifyProjectInstanceNew(d.Id(), call.param, client)
			},
			afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
				return err
			},
		}
	}
	return callback, err
}

func (s *EbsService) ModifyVolumeInfoCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"project_id":    {Ignore: true},
		"size":          {Ignore: true},
		"online_resize": {Ignore: true},
		"tags":          {Ignore: true},
	}
	req, err := SdkRequestAutoMapping(d, r, true, transform, nil, SdkReqParameter{
		false,
	})
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["VolumeId"] = d.Id()
		callback = ApiCall{
			param:  &req,
			action: "ModifyVolume",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.ebsconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyVolume(call.param)
				return resp, err
			},
			afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
				logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
				return err
			},
		}
	}
	return callback, err
}

func (s *EbsService) ModifyVolumeResizeCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {

	// 判断一下size是否变化，无变化就不调resize的接口
	// 这个显式判断如果不保留，在绑定状态下会触发一次resize调用，但是size没传就报错了
	if !d.HasChange("size") {
		return
	}
	transform := map[string]SdkReqTransform{
		"size": {},
		"online_resize": {
			forceUpdateParam: true,
		},
	}
	req, err := SdkRequestAutoMapping(d, r, true, transform, nil)
	var state interface{}
	state, err = s.checkVolumeState(d, "", []string{"available", "in-use"}, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return callback, err
	}
	if state.(map[string]interface{})["VolumeStatus"] == "available" {
		delete(req, "OnlineResize")
	}
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["VolumeId"] = d.Id()
		callback = ApiCall{
			param:  &req,
			action: "ResizeVolume",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.ebsconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ResizeVolume(call.param)
				return resp, err
			},
			afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
				logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
				_, err = s.checkVolumeState(d, "", []string{"available", "in-use"}, d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return err
				}
				return d.Set("online_resize", d.Get("online_resize"))
			},
		}
	}
	return callback, err
}

func (s *EbsService) ModifyVolume(d *schema.ResourceData, r *schema.Resource) (err error) {

	// a := d.HasChange("project_id")
	// b := d.HasChange("volume_desc")

	projectCall, err := s.ModifyVolumeProjectCall(d, r)
	if err != nil {
		return err
	}
	infoCall, err := s.ModifyVolumeInfoCall(d, r)
	if err != nil {
		return err
	}

	call, err := s.ModifyVolumeResizeCall(d, r)
	if err != nil {
		return err
	}

	tagsService := TagService{client: s.client}
	tagsCall, err := tagsService.ReplaceResourcesTagsWithResourceCall(d, r, "volume", true, false)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{projectCall, infoCall, call, tagsCall}, d, s.client, true)
}

func (s *EbsService) RemoveVolumeCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"VolumeId":    d.Id(),
		"ForceDelete": true,
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeleteVolume",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.ebsconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteVolume(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadVolume(d, "", false)
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading volume when delete %q, %s", d.Id(), callErr))
					}
				}
				_, callErr = call.executeCall(d, client, call)
				if callErr == nil {
					return nil
				}
				return resource.RetryableError(callErr)
			})
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}
	return callback, err
}

func (s *EbsService) RemoveVolume(d *schema.ResourceData) (err error) {
	call, err := s.RemoveVolumeCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

// start attach

func (s *EbsService) ReadVolumeAttach(d *schema.ResourceData, volumeId string, instanceId string) (data map[string]interface{}, err error) {
	data, err = s.ReadVolume(d, volumeId, false)
	if id, ok := data["InstanceId"]; ok {
		if id != instanceId {
			return data, fmt.Errorf("InstanceId %s not attach in Volume %s ", instanceId, volumeId)
		}
	} else {
		return data, fmt.Errorf("InstanceId %s not associate in Address %s ", instanceId, volumeId)
	}
	flag, err := getSdkValue("Attachment.0.DeleteWithInstance", data)
	if err != nil {
		return data, err
	}
	data["DeleteWithInstance"] = flag
	return data, err
}

func (s *EbsService) ReadAndSetVolumeAttach(d *schema.ResourceData, r *schema.Resource) (err error) {
	data, err := s.ReadVolumeAttach(d, d.Get("volume_id").(string), d.Get("instance_id").(string))
	if err != nil {
		return err
	}
	SdkResponseAutoResourceData(d, r, data, nil)
	return err
}

func (s *EbsService) CreateVolumeAttachCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"delete_with_instance": {ValueFunc: func(data *schema.ResourceData) (interface{}, bool) {
			return data.Get("delete_with_instance"), true
		}},
	}
	req, err := SdkRequestAutoMapping(d, r, false, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})
	if err != nil {
		return callback, err
	}
	callback = ApiCall{
		param:  &req,
		action: "AttachVolume",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.ebsconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.AttachVolume(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				// errMessage := strings.ToLower(baseErr.Error())
				if strings.Contains(baseErr.Error(), "OperationFailedWithTradeInstanceError") {
					_, callErr := call.executeCall(d, client, call)
					if callErr == nil {
						return nil
					}
					return resource.RetryableError(callErr)
				} else {
					return resource.NonRetryableError(baseErr)
				}
			})
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			d.SetId(d.Get("volume_id").(string) + ":" + d.Get("instance_id").(string))
			_, err = s.checkVolumeState(d, d.Get("volume_id").(string), []string{"in-use"}, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return err
			}
			return err
		},
	}
	return callback, err
}

func (s *EbsService) CreateVolumeAttach(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.CreateVolumeAttachCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *EbsService) ModifyVolumeAttachInfoCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, true, nil, nil)
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["VolumeId"] = d.Get("volume_id")
		callback = ApiCall{
			param:  &req,
			action: "ModifyVolume",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.ebsconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyVolume(call.param)
				return resp, err
			},
			afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
				logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
				return err
			},
		}
	}
	return callback, err
}

func (s *EbsService) ModifyVolumeAttach(d *schema.ResourceData, r *schema.Resource) (err error) {
	infoCall, err := s.ModifyVolumeAttachInfoCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{infoCall}, d, s.client, true)
}

func (s *EbsService) RemoveVolumeAttachCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"VolumeId":   d.Get("volume_id"),
		"InstanceId": d.Get("instance_id"),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DetachVolume",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.ebsconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DetachVolume(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadVolumeAttach(d, d.Get("volume_id").(string), d.Get("instance_id").(string))
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading volume attach when delete %q, %s", d.Id(), callErr))
					}
				}
				_, callErr = call.executeCall(d, client, call)
				if callErr == nil {
					return nil
				}
				return resource.RetryableError(callErr)
			})
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			_, err = s.checkVolumeState(d, d.Get("volume_id").(string), []string{"available"}, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return err
			}
			return err
		},
	}
	return callback, err
}

func (s *EbsService) RemoveVolumeAttach(d *schema.ResourceData) (err error) {
	call, err := s.RemoveVolumeAttachCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}
