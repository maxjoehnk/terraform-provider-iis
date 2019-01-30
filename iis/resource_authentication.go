package iis

import (
	"github.com/hashicorp/terraform/helper/schema"
	iis "github.com/maxjoehnk/microsoft-iis-administration"
	"log"
)

func resourceAuthentication() *schema.Resource {
	return &schema.Resource{
		Create: resourceAuthenticationCreate,
		Read:   resourceAuthenticationRead,
		Update: resourceAuthenticationUpdate,
		Delete: resourceAuthenticationDelete,

		Schema: map[string]*schema.Schema{
			"application": {
				Type:     schema.TypeString,
				Required: true,
			},
			"anonymous": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Default:  false,
							Optional: true,
						},
						"user": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Optional: true,
			},
			"basic": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Default:  false,
							Optional: true,
						},
						"default_domain": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"realm": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Optional: true,
			},
			"windows": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Default:  false,
							Optional: true,
						},
						"providers": {
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
					},
				},
				Optional: true,
			},
		},
	}
}

func resourceAuthenticationCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*iis.Client)
	auth, err := client.ReadAuthenticationFromApplication(d.Get("application").(string))
	if err != nil {
		return err
	}
	if err := updateAuthProviders(d, client, auth); err != nil {
		return err
	}
	d.SetId(auth.ID)
	return nil
}

func resourceAuthenticationRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*iis.Client)
	auth, err := client.ReadAuthentication(d.Id())
	if err != nil {
		return err
	}
	if err = readAuthenticationProvider(d, "anonymous", buildAnonymousAuthProvider(client, &auth)); err != nil {
		return err
	}
	if err = readAuthenticationProvider(d, "basic", buildBasicAuthProvider(client, &auth)); err != nil {
		return err
	}
	if err = readAuthenticationProvider(d, "windows", buildWindowsAuthProvider(client, &auth)); err != nil {
		return err
	}
	d.SetId(auth.ID)
	return nil
}

func resourceAuthenticationUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*iis.Client)
	auth, err := client.ReadAuthentication(d.Id())
	if err != nil {
		return err
	}
	return updateAuthProviders(d, client, auth)
}

func resourceAuthenticationDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func updateAuthProviders(d *schema.ResourceData, client *iis.Client, auth iis.Authentication) error {
	if err := updateAnonymousAuthentication(d, client, auth); err != nil {
		return err
	}
	if err := updateBasicAuthentication(d, client, auth); err != nil {
		return err
	}
	if err := updateWindowsAuthentication(d, client, auth); err != nil {
		return err
	}
	return nil
}

func updateAnonymousAuthentication(d *schema.ResourceData, client *iis.Client, auth iis.Authentication) error {
	if !d.HasChange("anonymous") {
		log.Println("no changes for anonymous authentication")
		return nil
	}
	anonymous, err := client.ReadAnonymousAuthentication(&auth)
	if err != nil {
		return err
	}
	anonymousMap := getNestedMap(d, "anonymous")
	anonymous.Enabled = anonymousMap["enabled"].(bool)
	anonymous.User = anonymousMap["user"].(string)

	_, err = client.UpdateAnonymousAuthentication(&anonymous)

	return err
}

func updateBasicAuthentication(d *schema.ResourceData, client *iis.Client, auth iis.Authentication) error {
	if !d.HasChange("basic") {
		log.Println("no changes for basic authentication")
		return nil
	}
	basic, err := client.ReadBasicAuthentication(&auth)
	if err != nil {
		return err
	}
	basicMap := getNestedMap(d, "basic")
	basic.Enabled = basicMap["enabled"].(bool)
	basic.DefaultLogonDomain = basicMap["default_domain"].(string)
	basic.Realm = basicMap["realm"].(string)

	_, err = client.UpdateBasicAuthentication(&basic)

	return err
}

func updateWindowsAuthentication(d *schema.ResourceData, client *iis.Client, auth iis.Authentication) error {
	if !d.HasChange("windows") {
		log.Println("no changes for windows authentication")
		return nil
	}
	windows, err := client.ReadWindowsAuthentication(&auth)
	if err != nil {
		return err
	}
	windowsMap := getNestedMap(d, "windows")
	windows.Enabled = windowsMap["enabled"].(bool)

	_, err = client.UpdateWindowsAuthentication(&windows)

	return err
}

func buildAnonymousAuthProvider(client *iis.Client, auth *iis.Authentication) FetchAuthProvider {
	return func() (AuthProvider, error) {
		return client.ReadAnonymousAuthentication(auth)
	}
}

func buildBasicAuthProvider(client *iis.Client, auth *iis.Authentication) FetchAuthProvider {
	return func() (AuthProvider, error) {
		return client.ReadBasicAuthentication(auth)
	}
}

func buildWindowsAuthProvider(client *iis.Client, auth *iis.Authentication) FetchAuthProvider {
	return func() (AuthProvider, error) {
		return client.ReadWindowsAuthentication(auth)
	}
}

type AuthProvider interface {
	ToMap() map[string]interface{}
}
type FetchAuthProvider func() (AuthProvider, error)

func readAuthenticationProvider(d *schema.ResourceData, key string, fetch FetchAuthProvider) error {
	provider, err := fetch()
	if err != nil {
		return err
	}
	providerList := make([]map[string]interface{}, 1)
	providerList = append(providerList, provider.ToMap())
	if err := d.Set(key, providerList); err != nil {
		return err
	}
	return nil
}
