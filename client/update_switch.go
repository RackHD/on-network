/* 
 * on-network api
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * OpenAPI spec version: 0.0.1
 * 
 * Generated by: https://github.com/swagger-api/swagger-codegen.git
 */

package swagger

type UpdateSwitch struct {

	Endpoint SwitchEndpoint `json:"endpoint"`

	ImageURL string `json:"imageURL"`

	SwitchModel string `json:"switchModel"`
}
