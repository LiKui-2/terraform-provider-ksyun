/*
Provides a Load Balancer Listener resource.

# Example Usage

```hcl

	resource "ksyun_certificate" "default" {
	  certificate_name="kkk"
	  public_key="-----BEGIN CERTIFICATE-----\nMIIE9zCCA9+gAwIBAgIQOJzS+B180J8Fyp3N2EQwTDANBgkqhkiG9w0BAQsFADBS\nMQswCQYDVQQGEwJDTjEaMBgGA1UEChMRV29TaWduIENBIExpbWl0ZWQxJzAlBgNV\nBAMTHldvU2lnbiBDbGFzcyAzIE9WIFNlcnZlciBDQSBHMjAeFw0xNTEyMzExMDA3\nMTlaFw0xOTAzMzExMDA3MTlaMHYxCzAJBgNVBAYTAkNOMRAwDgYDVQQIDAdUaWFu\namluMRAwDgYDVQQHDAdUaWFuamluMSswKQYDVQQKDCJUaWFuamluIFN1aXl1ZSBU\nZWNobm9sb2d5IENvLixMdGQuMRYwFAYDVQQDDA0qLnRpc2dhbWUuY29tMIIBIjAN\nBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0Gjl8buPjyFbLXNI2ie07gVHRGEv\nKbE8+wqVS/Uyi0AS0LqK+h37rHi1USizD8GTY2NNh6KbemfgflhiuxAsXTAtDzmB\nGkD8Auws68tVlu+ur1uht1gYtnTYldhi5c6EmOotTB0E4YtMQbYeTAqKGeYVDO00\nIF5scI3eVDQgw/qsJfOoUkjcM9VfYyalarkWo2A4tLrR527qkBtYmApLaHYY7Zmd\nQlV39bUktG8Pgbmvi+ycFfjhpACtGcoJKEfsydWEjEklQQDxRe46cb0Jkg2cpJ4J\nEF1YDIdh3AAsNgYEE7MdVhhYEuKgy5DqTtuPPTOVjh9fMtWo/u9a9VhPjwIDAQAB\no4IBozCCAZ8wCwYDVR0PBAQDAgWgMB0GA1UdJQQWMBQGCCsGAQUFBwMCBggrBgEF\nBQcDATAJBgNVHRMEAjAAMB0GA1UdDgQWBBQw/Pm54BOxQMFwzJOeiaZnXRKdRjAf\nBgNVHSMEGDAWgBT5i+wEOGo/qgbGlK1zlSqwyOa4+zBzBggrBgEFBQcBAQRnMGUw\nLwYIKwYBBQUHMAGGI2h0dHA6Ly9vY3NwMS53b3NpZ24uY29tL2NhNi9zZXJ2ZXIz\nMDIGCCsGAQUFBzAChiZodHRwOi8vYWlhMS53b3NpZ24uY29tL2NhNi5zZXJ2ZXIz\nLmNlcjA4BgNVHR8EMTAvMC2gK6AphidodHRwOi8vY3JsczEud29zaWduLmNvbS9j\nYTYtc2VydmVyMy5jcmwwJQYDVR0RBB4wHIINKi50aXNnYW1lLmNvbYILdGlzZ2Ft\nZS5jb20wUAYDVR0gBEkwRzAIBgZngQwBAgIwOwYMKwYBBAGCm1EGAwIBMCswKQYI\nKwYBBQUHAgEWHWh0dHA6Ly93d3cud29zaWduLmNvbS9wb2xpY3kvMA0GCSqGSIb3\nDQEBCwUAA4IBAQB5jIzf1Q4+IK+A+iicyznJn4kl56TMu8F2++zhWAwUP3ZyzJr3\nZaVkcfN+P5zRCCwy40+HHUb+zxQc8NTYLl88IBGyO3asaKZRzGlI8TkIXkEY2tlf\nFCZfAOJIwITwqNuepMlTyOjuqxhwzyr9Z2GASJ7Coqtrj6l6OoHvBNS9vNWziP1J\ngJ/cDpV4z02SY/fVw4udlT5J6FTGIOmMucnlh8CGsN6oFCPItIjVZhLGwgZbyNrz\nP6/4rdVZ2fVk8Q5Hn5arTKcwIOsroNxxPxLMxV5DNFwtJZ4gxcYz0o75VY/X9VYW\nWYdRxC4CjnSn/uVleWJBFcR0gj6vBPTWhQ4V\n-----END CERTIFICATE-----\n\n-----BEGIN CERTIFICATE-----\nMIIFozCCA4ugAwIBAgIQdZbCPvqJWUVuefcXus9k8zANBgkqhkiG9w0BAQsFADBV\nMQswCQYDVQQGEwJDTjEaMBgGA1UEChMRV29TaWduIENBIExpbWl0ZWQxKjAoBgNV\nBAMTIUNlcnRpZmljYXRpb24gQXV0aG9yaXR5IG9mIFdvU2lnbjAeFw0xNDExMDgw\nMDU4NThaFw0yOTExMDgwMDU4NThaMFIxCzAJBgNVBAYTAkNOMRowGAYDVQQKExFX\nb1NpZ24gQ0EgTGltaXRlZDEnMCUGA1UEAxMeV29TaWduIENsYXNzIDMgT1YgU2Vy\ndmVyIENBIEcyMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1nSHr5nA\nV5aZwol0PJJVmb8fBwA1BSaWFlsDwUI3M74/DU//u5QmkdcUFngb9xOiS0zlXKcQ\nQDVZMNF3meOdKcK+MZW9kmFbsCP7Z1jVUuR7L/BzHHOUVbrIaFkCEBDk9xHww7bX\nrlaAAJ5lZKaDkUHm7ad6ZaUfMC4TPL/fY5fzlvBSMrT0e5hX7TZP9yFKKJ3dHJKz\nTY2cWIsXIdjcobeuc3iKxLbpfyiOmtUunjnp2ll048iXEDKUGVnUD4lXROblKxcw\nYlKYf6sNpQHqBEHK+hMOO4cGur1HMddjAwH0vqE3EZ8eAZVODz9UHpKmnzCM/pjo\nVpZmBOE1/lmsVwIDAQABo4IBcDCCAWwwDgYDVR0PAQH/BAQDAgEGMB0GA1UdJQQW\nMBQGCCsGAQUFBwMCBggrBgEFBQcDATASBgNVHRMBAf8ECDAGAQH/AgEAMDAGA1Ud\nHwQpMCcwJaAjoCGGH2h0dHA6Ly9jcmxzMS53b3NpZ24uY29tL2NhMS5jcmwwbQYI\nKwYBBQUHAQEEYTBfMCcGCCsGAQUFBzABhhtodHRwOi8vb2NzcDEud29zaWduLmNv\nbS9jYTEwNAYIKwYBBQUHMAKGKGh0dHA6Ly9haWExLndvc2lnbi5jb20vY2ExZzIt\nc2VydmVyMy5jZXIwHQYDVR0OBBYEFPmL7AQ4aj+qBsaUrXOVKrDI5rj7MB8GA1Ud\nIwQYMBaAFOFmzw7R8bNLtwYgFP6HEtX2/vs+MEYGA1UdIAQ/MD0wOwYMKwYBBAGC\nm1EGAwIBMCswKQYIKwYBBQUHAgEWHWh0dHA6Ly93d3cud29zaWduLmNvbS9wb2xp\nY3kvMA0GCSqGSIb3DQEBCwUAA4ICAQBeZ7p4MgW2t6/n3mp6gmQOoAvynpq6xitv\nVjq0YlerfK1gUJY0nKOIz9mPUK/28AA2Gx8fh1U8YJrwsA2agC2KO74Fs9eggLa4\nGetR2+xkVPEaiUpIoU0/MX3EeZRL8d6rg69fhr6WHLM+HOe8lrLoWqy1WMs8Vm8K\np6XQNomCJoy5H7brj354/FuLeRzW30enVvSYTsep1Q51VgZ/tDdGCMbpT4tbQxzg\nRT6VIHHAHJgW7/J436xNu79WDs+Fr8+/BO1ya/0fVw5YkUQRWDtiOwl4s6R1auyz\nwisyzLONw6Nu3IrV6ErEC3vbMF2VM8PRo2lkW6iqlkhzc+PJuSTfF3Wqrwc6z76b\nioCnv3zi6Srm/bAs5+bmfrM1FWUA9OE5cw4oS/AMmJ466857ep5AwVBllprnS3fN\n3ct9l7TqCbLpSSjDMOCHFfAm6tgD/ezaCINl3HfFbj0094fDHB0mM+wzrMaZU6tg\n9LDZ7mRaMwdwE3SIB/WG+RjTskfIrgNKU94cZdYKLjpRk+63428K++n+Tui7HcKX\nqwq57TYyG02hzAOmnbPZHNVn4o90PJIqdLFWUN9TFdch1uvz+2PjICwKdDcLwaE1\naoRw9EX4sraBSar9VEWQTecEB194FN06uyv5clDsaOo8qNGAu741Q5fDMrL1qq3J\nf4OffWkeFQ==\n-----END CERTIFICATE-----\n\n-----BEGIN CERTIFICATE-----\nMIIGXDCCBESgAwIBAgIHGcKFMOk7NjANBgkqhkiG9w0BAQsFADB9MQswCQYDVQQG\nEwJJTDEWMBQGA1UEChMNU3RhcnRDb20gTHRkLjErMCkGA1UECxMiU2VjdXJlIERp\nZ2l0YWwgQ2VydGlmaWNhdGUgU2lnbmluZzEpMCcGA1UEAxMgU3RhcnRDb20gQ2Vy\ndGlmaWNhdGlvbiBBdXRob3JpdHkwHhcNMDYwOTE3MjI0NjM2WhcNMTkxMjMxMjM1\nOTU5WjBVMQswCQYDVQQGEwJDTjEaMBgGA1UEChMRV29TaWduIENBIExpbWl0ZWQx\nKjAoBgNVBAMTIUNlcnRpZmljYXRpb24gQXV0aG9yaXR5IG9mIFdvU2lnbjCCAiIw\nDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBAL3Kjay4kRVWl3trXHrC3mvZobDD\nECP6p6GyzDH6PtmmKW8WPeBr+LhAX9s5qAB6i6BNVH3CInj8jgm4qIXXzJWXS3TY\nnn7wAOQOia5JKEQaEJkyDyWIU6QNsw8SCBYLA3EnHH/h29L9Z2jEBV0KDl1w19iX\noLxTQZqRjfSeNmZ6flbBkF/msWggNqSMJCwsRwtZdmYwtb7e7Y/4ndO7ATDm8vMO\n4CySgPOF+SiKtFQumu33dvwVaBbrSmzrLhKP1M/+DMdcHQt+BTK+XrAJKkLVyU6Q\ns1kNu3p+zdUIWrR/2BxpEfknD3sGr1SDGHvh3VR6UWhud/zGv1JKZkahsmcau6NP\nd6C+Xf/8VgtDcneQyp758jn1Dan06tfnsxAvMEI3IcwwcMmGmA/MWE2Du33lGqU3\njbasMpcAOmNxJB6eN8T/dNQ3wOL+iEZgEd0IP1A2q7h6pJViam6wymohWmnz8/sd\ncDmV86dupoGJoYjFO3HKo1Lug7v9oHf05G/nQtttSpmKNEi8F9zkgAgitvIxwD8E\nPuufIHnWuAZkZAIx16nNUvuERWkJACrcVYvEBkZLwEodCVs5KP2pq84A+S5ISybm\nMEylWMq0RIJP55EeM8Owk/8R/IHSyh9xKd12T5Ilrx2Btw8vjMMGzC8no0rkDpm6\nfB5FH3+qGUWW/fw9AgMBAAGjggEHMIIBAzASBgNVHRMBAf8ECDAGAQH/AgECMA4G\nA1UdDwEB/wQEAwIBBjAdBgNVHQ4EFgQU4WbPDtHxs0u3BiAU/ocS1fb++z4wHwYD\nVR0jBBgwFoAUTgvvGqRAW6UXaYcwyjRoQ9BBrvIwaQYIKwYBBQUHAQEEXTBbMCcG\nCCsGAQUFBzABhhtodHRwOi8vb2NzcC5zdGFydHNzbC5jb20vY2EwMAYIKwYBBQUH\nMAKGJGh0dHA6Ly9haWEuc3RhcnRzc2wuY29tL2NlcnRzL2NhLmNydDAyBgNVHR8E\nKzApMCegJaAjhiFodHRwOi8vY3JsLnN0YXJ0c3NsLmNvbS9zZnNjYS5jcmwwDQYJ\nKoZIhvcNAQELBQADggIBALZt+HD74g1MmLMHSRX1BMRsysr1aKAI/hJtnAQGya2a\nkVI+eMRc7p9UHe7j8V4wyUnhOeCmnTZsV/rmNE9V6IeoLN0F8VgSkejKzih4j98H\nhQGl3EWWBdSAsisFmsuapYvgOmfmc0e+Sv0nsYjv5srPjQ4mn/pfV3itbf6umzUI\nscO6wQBKS30Uvffx01UYrNAzcIhtxAlxFKYrT4iB5wsAN6kVfX7XAZY/L697Yq4K\nSr9LOS41EIv+BDnkPDoMCVZAOrX0wmgMtflSze6d+Jj8eOdYR48cc1hpM6v/3d+O\nJAF3mBk6sGZ5vOEIow5PwQSz8wHI69NZHDXSkx5wZYJ/28/7yJkSYMNEbzqAS9e+\nIaoUemTL3TdDRVsyLkXw2VkfaxjwfOlVNhlhX7V98Y29iOR1S5jdJ7DkhEQqYYRX\nBYIRH6o1WPMgDq9Z7/pVcnINJtCbU0mszjcuZWH/9uwb6vbxptPRtXu+NfQiwbyN\nAb1oXoMNL+zW2mMMJ9FUPuSo085LMriRlP/7W0ktdRiounGaO67ZwKlPh5Hti3tr\nIJiJOYNPgMRpzBfJyE6+5KmlgXZwBgQyzYNl9Lx9PhO80uhvY6q1O9qNhjKCeJ3Z\nzP+/V2R07Sg9RGIVYUv3lLANKmcc8MubpZK/+EFawT1g7Z+7uG2bzqlqFj9+6gbx\n-----END CERTIFICATE-----"
	  private_key="-----BEGIN RSA PRIVATE KEY-----\nMIIEpQIBAAKCAQEA0Gjl8buPjyFbLXNI2ie07gVHRGEvKbE8+wqVS/Uyi0AS0LqK\n+h37rHi1USizD8GTY2NNh6KbemfgflhiuxAsXTAtDzmBGkD8Auws68tVlu+ur1uh\nt1gYtnTYldhi5c6EmOotTB0E4YtMQbYeTAqKGeYVDO00IF5scI3eVDQgw/qsJfOo\nUkjcM9VfYyalarkWo2A4tLrR527qkBtYmApLaHYY7ZmdQlV39bUktG8Pgbmvi+yc\nFfjhpACtGcoJKEfsydWEjEklQQDxRe46cb0Jkg2cpJ4JEF1YDIdh3AAsNgYEE7Md\nVhhYEuKgy5DqTtuPPTOVjh9fMtWo/u9a9VhPjwIDAQABAoIBAB84t8oBCT6uBpTG\ngpF8BVTrIYQS8xfR34rUSQ16DDpkPO1cJNnXiA6Mfs1LxZeErmnjRIgfSh+KVE1Z\nPGE9yQek4fs36ClgtuBRNWGVsprMCWVn2tA2uG0NEC1Sn25CBAdMfcC7UDHhxmjJ\nubw2j/3VRC4NkxOOUbyC5F3E+WuiasL5NcH/gWrT9jdHIcmpJXhtw65gJOt3HEkf\n9RcQYjCSS8WVMtS7xRhK+IeXK+Vbn3pcXQ9Q6icce9+LSboeyo0GPi4GC3KcA3/R\n5t5Zwl4Yo04J2NfjJbSquPZnvdPyaS2iwvtatZUuUrYAiZKwjhfDOmD6o6fDr1wL\nBBSCJUECgYEA8iGUflCTQIoCYl+c6xhXsf8YUAKd/UIcRZVOKdna3rV84eopZYNS\nNiw139yyJHUeGtzrMDvtCiR7btxyyzYFoK/7NhMzJE8zia/UG3zrglNOJt2VwhsU\nwRaSl4tUNOAvf8BlVPW5WXrqCk5aD6OodnzK+DmYXE0moeoJHAMiRv8CgYEA3Fja\nolXYWsf7iHr5+iGzckODGGiWQI1wECPWJYm8qfMNVeOMyOD/jxII688+1cI0GFxN\nGtIEFd4njsHSR96934gF7IeujbrnRIvlfAkxQV57AP9602cSuvLQ696nycPdydbF\n+GcbOnazzlUUPP1vdwdmu8elZdBEgoaTHfK6B3ECgYEA05Wegmb81lgDT12H7TR/\nZY6p+zjeQHJl7DRVcmLqTNVBRNVvyUJhM++cQHxFu3AQl57XcnXbZJKOvkirk4Io\nlstRdWZ/uUnwmm/opQCbeG49i970QAOUNkr4XK8nLXF+cF967SwxBM1Q+SKQtrvn\nuWrBnvoNdxMAIFs4DJ98c5cCgYEAm82oKokQxoAJd2OdRiR2QiFCnQu20kYwKvlr\n/nb6FCFsGIMhlRijG4LhE/wirfr34xHA42oEwYGn7uVVzsPM9jW0Gp+F6WlzBaD1\nz1KgpVwtFXOQYdLMB6yR1XZGpf/83y8iJJajRh+Q3CCEguug2UU+eyCb7vXou3J0\nrARpTcECgYEAzGugegAYtrlF+c+zFGySzgtIxYld6H+T9cDgrUtWAOp+P3SwbvLl\n5PdRJ183xZo6s4O+Ptv5gl+XcNFIM/xiSmuyogFcwNBifr45anqmOlok/Y0N1cya\nEXN6Umpw6rO0b9aCuhKAJES+QYTZ4jzn5NYphC2t4yv+0KFlUHJWEok=\n-----END RSA PRIVATE KEY-----"
	}

	resource "ksyun_lb_listener" "default" {
	  listener_name = "tf-xun"
	  listener_port = "8000"
	  listener_protocol = "HTTPS"
	  listener_state = "stop"
	  load_balancer_id = "7fae85e4-ab1a-415c-aef9-03a402c79d97",
	  method = "RoundRobin"
	  certificate_id = "${ksyun_certificate.default.id}"
	  session {
	    session_state = "stop"
	    session_persistence_period = 100
	    cookie_type = "ImplantCookie"
	    cookie_name = "cookiexunqq"
	  }
	}

```

# Import

LB Listener can be imported using the `id`, e.g.

```
$ terraform import ksyun_lb_listener.example vserver-abcdefg
```
*/
package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunListener() *schema.Resource {
	entry := resourceKsyunHealthCheck().Schema
	for k, v := range entry {
		if k == "listener_id" {
			v.Required = false
			v.Computed = true
		} else {
			v.ForceNew = false
			v.DiffSuppressFunc = nil
		}
		switch k {
		case "lb_type":
			delete(entry, k)
		}
	}
	return &schema.Resource{
		Create: resourceKsyunListenerCreate,
		Read:   resourceKsyunListenerRead,
		Update: resourceKsyunListenerUpdate,
		Delete: resourceKsyunListenerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the LB.",
			},
			"listener_state": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "start",
				ValidateFunc: validation.StringInSlice([]string{
					"start",
					"stop",
				}, false),
				Description: "The state of listener.Valid Values:'start', 'stop'.",
			},
			"listener_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of listener.",
			},
			"listener_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "TCP",
				ValidateFunc: validation.StringInSlice([]string{
					"TCP",
					"UDP",
					"HTTP",
					"HTTPS",
				}, false),
				ForceNew:    true,
				Description: "The protocol of listener.Valid Values:'TCP', 'UDP', 'HTTP', 'HTTPS'.",
			},
			"certificate_id": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: lbListenerDiffSuppressFunc,
				Description:      "The ID of certificate.",
			},
			"listener_port": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
				Description:  "The protocol port of listener.",
			},
			"method": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "RoundRobin",
				// ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"RoundRobin",
					"LeastConnections",
					"MasterSlave",
					"QUIC_CID",
				}, false),
				Description: "Forwarding mode of listener.Valid Values:'RoundRobin', 'LeastConnections', 'MasterSlave', 'QUIC_CID'.",
			},

			"enable_http2": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          true,
				DiffSuppressFunc: lbListenerDiffSuppressFunc,
				Description:      "whether enable to HTTP2.",
			},

			"tls_cipher_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "TlsCipherPolicy1.0",
				ValidateFunc: validation.StringInSlice([]string{
					"TlsCipherPolicy1.0",
					"TlsCipherPolicy1.1",
					"TlsCipherPolicy1.2",
					"TlsCipherPolicy1.2-strict",
				}, false),
				DiffSuppressFunc: lbListenerDiffSuppressFunc,
				Description:      "TLS cipher policy, valid values:'TlsCipherPolicy1.0','TlsCipherPolicy1.1','TlsCipherPolicy1.2','TlsCipherPolicy1.2-strict'.",
			},
			"http_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "HTTP1.1",
				ValidateFunc: validation.StringInSlice([]string{
					"HTTP1.0",
					"HTTP1.1",
				}, false),
				DiffSuppressFunc: lbListenerDiffSuppressFunc,
				Description:      "HTTP protocol, valid values:'HTTP1.0','HTTP1.1'.",
			},

			"redirect_listener_id": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: lbListenerDiffSuppressFunc,
				Description:      "The ID of the redirect listener.",
			},

			"session": {
				Type:        schema.TypeList,
				MaxItems:    1,
				MinItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "session.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"session_state": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "stop",
							ValidateFunc: validation.StringInSlice([]string{
								"start",
								"stop",
							}, false),
							Description: "The state of session.Valid Values:'start', 'stop'.",
						},
						"session_persistence_period": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      3600,
							ValidateFunc: validation.IntBetween(1, 86400),
							Description:  "Session hold timeout.Valid Values:1-86400.",
						},
						"cookie_type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "ImplantCookie",
							ValidateFunc: validation.StringInSlice([]string{
								"ImplantCookie",
								"RewriteCookie",
							}, false),
							Description: "The type of cookie, valid values: 'ImplantCookie','RewriteCookie'.",
						},
						"cookie_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The name of cookie.",
						},
					},
				},
				DiffSuppressFunc: lbListenerDiffSuppressFunc,
			},

			"health_check": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: entry,
				},
				Deprecated:       "This parameter is deprecated and will be removed in a future version, use `ksyun_healthcheck` instead.",
				DiffSuppressFunc: lbListenerDiffSuppressFunc,
				Description:      "Health check.",
			},

			"load_balancer_acl_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of LB ACL.",
				// 设置optional+computed造成通过这个值可以绑定和修改，但是不能解绑，因此不设置computed
				// 但是需要注意，用ksyun_lb_listener_associate_acl会导致这个值必须设置，否则plan的时候会提示change
				// Computed: true,
			},

			"listener_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of listener.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time.",
			},
		},
	}
}

func resourceKsyunListenerCreate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.CreateListener(d, resourceKsyunListener())
	if err != nil {
		return fmt.Errorf("error on creating listener %q, %s", d.Id(), err)
	}
	return resourceKsyunListenerRead(d, meta)
}

func resourceKsyunListenerRead(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ReadAndSetListener(d, resourceKsyunListener())
	if err != nil {
		return fmt.Errorf("error on reading listener %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunListenerUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ModifyListener(d, resourceKsyunListener())
	if err != nil {
		return fmt.Errorf("error on updating listener %q, %s", d.Id(), err)
	}
	return resourceKsyunListenerRead(d, meta)
}

func resourceKsyunListenerDelete(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.RemoveListener(d)
	if err != nil {
		return fmt.Errorf("error on deleting listener %q, %s", d.Id(), err)
	}
	return err
}
