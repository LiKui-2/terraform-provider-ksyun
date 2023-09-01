package ksyun

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

type AlbListenerService struct {
	client *KsyunClient
}

func (s *AlbListenerService) createListenerCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"enable_http2": {
			ValueFunc: func(data *schema.ResourceData) (interface{}, bool) {
				return data.Get("enable_http2"), true
			},
		},
		"session": {
			Type: TransformListUnique,
		},
	}
	req, err := SdkRequestAutoMapping(d, r, false, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})

	if req["listener_protocol"] != "HTTPS" {
		delete(req, "EnableHttp2")
	}
	for k, v := range req {
		if strings.HasPrefix(k, "Session.") {
			req[strings.Replace(k, "Session.", "", -1)] = v
			delete(req, k)
		}
	}
	// if session is zero need set default SessionState stop
	if _, ok := req["SessionState"]; !ok {
		req["SessionState"] = "stop"
	}

	if err != nil {
		return callback, err
	}
	callback = ApiCall{
		param:  &req,
		action: "CreateAlbListener",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateAlbListener(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("AlbListener.AlbListenerId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			return d.Set("alb_listener_id", d.Id())
		},
	}
	return
}
func (s *AlbListenerService) CreateListener(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.createListenerCall(d, r)
	if err != nil {
		return err
	}
	err = ksyunApiCallNew([]ApiCall{
		call,
	}, d, s.client, false)
	return
}

func (s *AlbListenerService) readListeners(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)

	return pageQuery(condition, "MaxResults", "NextToken", 200, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		conn := s.client.slbconn
		action := "DescribeAlbListeners"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = conn.DescribeAlbListeners(nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = conn.DescribeAlbListeners(&condition)
			if err != nil {
				return data, err
			}
		}

		results, err = getSdkValue("AlbListenerSet", *resp)
		if err != nil {
			return data, err
		}
		data = results.([]interface{})
		return data, err
	})
}

func (s *AlbListenerService) readListener(d *schema.ResourceData, listenerId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if listenerId == "" {
		listenerId = d.Id()
	}
	req := map[string]interface{}{
		"AlbListenerId.1": listenerId,
	}
	results, err = s.readListeners(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("ALb listener %s not exist ", listenerId)
	}
	return
}

func (s *AlbListenerService) ReadAndSetListener(d *schema.ResourceData, r *schema.Resource) (err error) {
	data, err := s.readListener(d, "")
	if err != nil {
		return err
	}
	extra := map[string]SdkResponseMapping{
		"Session": {
			Field: "session",
			FieldRespFunc: func(i interface{}) interface{} {
				return []interface{}{
					i,
				}
			},
		},
	}
	SdkResponseAutoResourceData(d, r, data, extra)
	return
}

func (s *AlbListenerService) removeAlbListenerCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"AlbListenerId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeleteAlbListener",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteAlbListener(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.readListener(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading lb when delete %q, %s", d.Id(), callErr))
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
	return
}
func (s *AlbListenerService) RemoveListener(d *schema.ResourceData) (err error) {
	call, err := s.removeAlbListenerCall(d)
	if err != nil {
		return err
	}
	err = ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
	return
}

func (s *AlbListenerService) modifyListenerCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"session": {
			Type: TransformListUnique,
		},
	}
	req, err := SdkRequestAutoMapping(d, r, true, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})
	logger.Debug(logger.ReqFormat, "modifyAlbListener", req)
	if err != nil {
		return callback, err
	}
	// 特殊处理下"Session."
	for k, v := range req {
		if strings.HasPrefix(k, "Session.") {
			req[strings.Replace(k, "Session.", "", -1)] = v
			delete(req, k)
		}
	}

	logger.Debug(logger.RespFormat, "ModifyListenerCall", req, d.Get("session.0.cookie_name"))
	// 设置rewriteCookie的时候，如果之前cookiename没改，需要手动传入这个值
	if req["CookieType"] == "RewriteCookie" {
		if _, ok := req["CookieName"]; !ok {
			req["CookieName"] = d.Get("session.0.cookie_name")
		}
	}
	// 如果需要改cookieName，必须传入cookietype
	if _, ok := req["CookieName"]; ok {
		if _, ok := req["CookieType"]; !ok {
			req["CookieType"] = d.Get("session.0.cookie_type")
		}
	}

	if len(req) > 0 {
		req["AlbListenerId"] = d.Id()
		callback = ApiCall{
			param:  &req,
			action: "ModifyAlbListener",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.slbconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyAlbListener(call.param)
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

func (s *AlbListenerService) ModifyListener(d *schema.ResourceData, r *schema.Resource) (err error) {
	callbacks := []ApiCall{}
	var call ApiCall
	call, err = s.modifyListenerCall(d, r)
	if err != nil {
		return
	}
	callbacks = append(callbacks, call)
	err = ksyunApiCallNew(callbacks, d, s.client, true)
	return
}

func (s *AlbListenerService) ReadAndSetAlbListeners(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "AlbListenerId",
			Type:    TransformWithN,
		},
		"alb_id": {
			mapping: "load-balancer-id",
			Type:    TransformWithFilter,
		},
		"acl_id": {
			mapping: "load-balancer-acl-id",
			Type:    TransformWithFilter,
		},
		"protocol": {
			mapping: "listener-protocol",
			Type:    TransformWithFilter,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.readListeners(req)
	if err != nil {
		return err
	}
	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		nameField:   "AlbListenerName",
		idFiled:     "AlbListenerId",
		targetField: "listeners",
		extra:       map[string]SdkResponseMapping{},
	})
}
