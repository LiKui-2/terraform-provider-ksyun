package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"net"
	"regexp"
)

var validateName = validation.StringMatch(
	regexp.MustCompile(`^[A-Za-z0-9\p{Han}-_]{1,63}$`),
	"expected value to be 1 - 63 characters and only support chinese, english, numbers, '-', '_'",
)

// validateCIDRNetworkAddress ensures that the string value is a valid CIDR that
// represents a network address - it adds an error otherwise
func validateCIDRNetworkAddress(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, ipnet, err := net.ParseCIDR(value)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid CIDR, got error parsing: %s", k, err))
		return
	}

	if ipnet == nil || value != ipnet.String() {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid network CIDR, expected %q, got %q",
			k, ipnet, value))
	}

	return
}

func validateIpAddress(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	res := net.ParseIP(value)

	if res == nil {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid ip address, got error parsing: %s", k, value))
	}

	return
}

func validateSubnetType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "Reserve" && value != "Normal" && value != "Physical" {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid subnet type, got error parsing: %s", k, value))
	}
	return
}

func validateLbState(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "start" && value != "stop" {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid loadbalancer state, got error parsing: %s", k, value))
	}
	return
}
func validateLbType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "public" && value != "internal" {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid loadbalancer type, got error parsing: %s", k, value))
	}
	return
}

func validateRouteType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "InternetGateway" && value != "Tunnel" && value != "Host" && value != "Peering" && value != "DirectConnect" && value != "Vpn" {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid route type, got error parsing: %s", k, value))
	}
	return
}

func validateNatType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "public" {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid nat type, got error parsing: %s", k, value))
	}
	return
}

func validateNatMode(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "Vpc" && value != "Subnet" {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid nat mode, got error parsing: %s", k, value))
	}
	return
}

func validateNatIpNumber(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 1 || value > 10 {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid nat ip number in 1-10 and control by quota system, got error parsing: %d", k, value))
	}
	return
}

func validateNatBandWidth(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 1 || value > 15000 {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid nat bandwidth in 1-15000 and control by quota system, got error parsing: %d", k, value))
	}
	return
}

func validateNatChargeType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "Monthly" && value != "Peak" && value != "Daily" && value != "TrafficMonthly" &&
		value != "DailyPaidByTransfer" && value != "PostPaidByAdvanced95Peak" {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid nat charge type and control by price system, got error parsing: %s", k, value))
	}
	return
}

//校验Ks3 Bucket name
/*
func validateKs3BucketName(value string) error {
	if (len(value) < 3) || (len(value) > 63) { //3~63字符之间
		return fmt.Errorf("%q must contain from 3 to 63 characters", value)
	}
	if !regexp.MustCompile(`^[0-9a-z-.]+$`).MatchString(value) { //小写和数字
		return fmt.Errorf("only lowercase alphanumeric characters and hyphens allowed in %q", value)
	}
	if regexp.MustCompile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`).MatchString(value) { //不能是IP
		return fmt.Errorf("%q must not be formatted as an IP address", value)
	}
	if strings.HasPrefix(value, `.`) { //不能以点开头
		return fmt.Errorf("%q cannot start with a period", value)
	}
	if strings.HasSuffix(value, `.`) { //不能以点结尾
		return fmt.Errorf("%q cannot end with a period", value)
	}
	if strings.Contains(value, `..`) { //不能包含两个点
		return fmt.Errorf("%q can be only one period between labels", value)
	}
	return nil
}

func validateKs3BucketLifecycleTransitionStorageClass() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		s3.TransitionStorageClassGlacier,
		s3.TransitionStorageClassStandardIa,
		s3.TransitionStorageClassOnezoneIa,
		s3.TransitionStorageClassIntelligentTiering,
		s3.TransitionStorageClassDeepArchive,
	}, false)
}
func validateKs3BucketLifecycleTimestamp(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00Z", value))
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q cannot be parsed as RFC3339 Timestamp Format", value))
	}

	return
}

*/
