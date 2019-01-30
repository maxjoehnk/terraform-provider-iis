package iis

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/maxjoehnk/microsoft-iis-administration"
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
	anonymousAuthProvider := buildAnonymousAuthProvider(client, &auth)
	basicAuthProvider := buildBasicAuthProvider(client, &auth)
	windowsAuthProvider := buildWindowsAuthProvider(client, &auth)

	if err := updateAuthenticationProvider(d, client, "anonymous", anonymousAuthProvider, updateAnonymousAuthentication); err != nil {
		return err
	}
	if err := updateAuthenticationProvider(d, client, "basic", basicAuthProvider, updateBasicAuthentication); err != nil {
		return err
	}
	if err := updateAuthenticationProvider(d, client, "windows", windowsAuthProvider, updateWindowsAuthentication); err != nil {
		return err
	}
	return nil
}

func updateAnonymousAuthentication(client *iis.Client, auth interface{}, data map[string]interface{}) error {
	anonymous := auth.(iis.AnonymousAuthentication)
	anonymous.Enabled = data["enabled"].(bool)
	anonymous.User = data["user"].(string)

	_, err := client.UpdateAnonymousAuthentication(&anonymous)

	return err
}

func updateBasicAuthentication(client *iis.Client, auth interface{}, data map[string]interface{}) error {
	basic := auth.(iis.BasicAuthentication)
	basic.Enabled = data["enabled"].(bool)
	basic.DefaultLogonDomain = data["default_domain"].(string)
	basic.Realm = data["realm"].(string)

	_, err := client.UpdateBasicAuthentication(&basic)

	return err
}

func updateWindowsAuthentication(client *iis.Client, auth interface{}, data map[string]interface{}) error {
	windows := auth.(iis.WindowsAuthentication)
	windows.Enabled = data["enabled"].(bool)

	_, err := client.UpdateWindowsAuthentication(&windows)

	return err
}

func updateAuthenticationProvider(d *schema.ResourceData, client *iis.Client, key string, fetch FetchAuthProvider, update UpdateAuthProvider) error {
	if !d.HasChange(key) {
		log.Printf("no changes for %s authentication", key)
		return nil
	}
	provider, err := fetch()
	if err != nil {
		return err
	}

	if !hasNestedMap(d, key) {
		return nil
	}

	return update(client, provider, getNestedMap(d, key))
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

type UpdateAuthProvider func(*iis.Client, interface{}, map[string]interface{}) error

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
