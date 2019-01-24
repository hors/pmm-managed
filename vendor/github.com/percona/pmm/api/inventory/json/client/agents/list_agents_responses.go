// Code generated by go-swagger; DO NOT EDIT.

package agents

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	strfmt "github.com/go-openapi/strfmt"
)

// ListAgentsReader is a Reader for the ListAgents structure.
type ListAgentsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListAgentsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewListAgentsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewListAgentsOK creates a ListAgentsOK with default headers values
func NewListAgentsOK() *ListAgentsOK {
	return &ListAgentsOK{}
}

/*ListAgentsOK handles this case with default header values.

(empty)
*/
type ListAgentsOK struct {
	Payload *ListAgentsOKBody
}

func (o *ListAgentsOK) Error() string {
	return fmt.Sprintf("[POST /v1/inventory/Agents/List][%d] listAgentsOK  %+v", 200, o.Payload)
}

func (o *ListAgentsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(ListAgentsOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*ExternalAgentItems0 ExternalAgent does not run on any Inventory Node.
swagger:model ExternalAgentItems0
*/
type ExternalAgentItems0 struct {

	// Unique randomly generated instance identifier.
	AgentID string `json:"agent_id,omitempty"`

	// URL for scraping metrics.
	MetricsURL int64 `json:"metrics_url,omitempty"`
}

// Validate validates this external agent items0
func (o *ExternalAgentItems0) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *ExternalAgentItems0) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *ExternalAgentItems0) UnmarshalBinary(b []byte) error {
	var res ExternalAgentItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*ListAgentsBody list agents body
swagger:model ListAgentsBody
*/
type ListAgentsBody struct {

	// Return only Agents that provide insights for that Node.
	NodeID string `json:"node_id,omitempty"`

	// Return only Agents running on that Node.
	RunsOnNodeID string `json:"runs_on_node_id,omitempty"`

	// Return only Agents that provide insights for that Service.
	ServiceID string `json:"service_id,omitempty"`
}

// Validate validates this list agents body
func (o *ListAgentsBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *ListAgentsBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *ListAgentsBody) UnmarshalBinary(b []byte) error {
	var res ListAgentsBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*ListAgentsOKBody list agents o k body
swagger:model ListAgentsOKBody
*/
type ListAgentsOKBody struct {

	// external agent
	ExternalAgent []*ExternalAgentItems0 `json:"external_agent"`

	// mysqld exporter
	MysqldExporter []*MysqldExporterItems0 `json:"mysqld_exporter"`

	// node exporter
	NodeExporter []*NodeExporterItems0 `json:"node_exporter"`

	// pmm agent
	PMMAgent []*PMMAgentItems0 `json:"pmm_agent"`

	// rds exporter
	RDSExporter []*RDSExporterItems0 `json:"rds_exporter"`
}

// Validate validates this list agents o k body
func (o *ListAgentsOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateExternalAgent(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateMysqldExporter(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateNodeExporter(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validatePMMAgent(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateRDSExporter(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *ListAgentsOKBody) validateExternalAgent(formats strfmt.Registry) error {

	if swag.IsZero(o.ExternalAgent) { // not required
		return nil
	}

	for i := 0; i < len(o.ExternalAgent); i++ {
		if swag.IsZero(o.ExternalAgent[i]) { // not required
			continue
		}

		if o.ExternalAgent[i] != nil {
			if err := o.ExternalAgent[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("listAgentsOK" + "." + "external_agent" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (o *ListAgentsOKBody) validateMysqldExporter(formats strfmt.Registry) error {

	if swag.IsZero(o.MysqldExporter) { // not required
		return nil
	}

	for i := 0; i < len(o.MysqldExporter); i++ {
		if swag.IsZero(o.MysqldExporter[i]) { // not required
			continue
		}

		if o.MysqldExporter[i] != nil {
			if err := o.MysqldExporter[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("listAgentsOK" + "." + "mysqld_exporter" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (o *ListAgentsOKBody) validateNodeExporter(formats strfmt.Registry) error {

	if swag.IsZero(o.NodeExporter) { // not required
		return nil
	}

	for i := 0; i < len(o.NodeExporter); i++ {
		if swag.IsZero(o.NodeExporter[i]) { // not required
			continue
		}

		if o.NodeExporter[i] != nil {
			if err := o.NodeExporter[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("listAgentsOK" + "." + "node_exporter" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (o *ListAgentsOKBody) validatePMMAgent(formats strfmt.Registry) error {

	if swag.IsZero(o.PMMAgent) { // not required
		return nil
	}

	for i := 0; i < len(o.PMMAgent); i++ {
		if swag.IsZero(o.PMMAgent[i]) { // not required
			continue
		}

		if o.PMMAgent[i] != nil {
			if err := o.PMMAgent[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("listAgentsOK" + "." + "pmm_agent" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (o *ListAgentsOKBody) validateRDSExporter(formats strfmt.Registry) error {

	if swag.IsZero(o.RDSExporter) { // not required
		return nil
	}

	for i := 0; i < len(o.RDSExporter); i++ {
		if swag.IsZero(o.RDSExporter[i]) { // not required
			continue
		}

		if o.RDSExporter[i] != nil {
			if err := o.RDSExporter[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("listAgentsOK" + "." + "rds_exporter" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *ListAgentsOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *ListAgentsOKBody) UnmarshalBinary(b []byte) error {
	var res ListAgentsOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*MysqldExporterItems0 MySQLdExporter runs on Generic or Container Node and exposes MySQL and AmazonRDSMySQL Service metrics.
swagger:model MysqldExporterItems0
*/
type MysqldExporterItems0 struct {

	// Unique randomly generated instance identifier.
	AgentID string `json:"agent_id,omitempty"`

	// Listen port for scraping metrics.
	ListenPort int64 `json:"listen_port,omitempty"`

	// MySQL password for scraping metrics.
	Password string `json:"password,omitempty"`

	// Node identifier where this instance runs.
	RunsOnNodeID string `json:"runs_on_node_id,omitempty"`

	// Service identifier.
	ServiceID string `json:"service_id,omitempty"`

	// AgentStatus represents actual Agent process status.
	// Enum: [AGENT_STATUS_INVALID STARTING RUNNING WAITING STOPPING DONE]
	Status *string `json:"status,omitempty"`

	// MySQL username for scraping metrics.
	Username string `json:"username,omitempty"`
}

// Validate validates this mysqld exporter items0
func (o *MysqldExporterItems0) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateStatus(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var mysqldExporterItems0TypeStatusPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["AGENT_STATUS_INVALID","STARTING","RUNNING","WAITING","STOPPING","DONE"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		mysqldExporterItems0TypeStatusPropEnum = append(mysqldExporterItems0TypeStatusPropEnum, v)
	}
}

const (

	// MysqldExporterItems0StatusAGENTSTATUSINVALID captures enum value "AGENT_STATUS_INVALID"
	MysqldExporterItems0StatusAGENTSTATUSINVALID string = "AGENT_STATUS_INVALID"

	// MysqldExporterItems0StatusSTARTING captures enum value "STARTING"
	MysqldExporterItems0StatusSTARTING string = "STARTING"

	// MysqldExporterItems0StatusRUNNING captures enum value "RUNNING"
	MysqldExporterItems0StatusRUNNING string = "RUNNING"

	// MysqldExporterItems0StatusWAITING captures enum value "WAITING"
	MysqldExporterItems0StatusWAITING string = "WAITING"

	// MysqldExporterItems0StatusSTOPPING captures enum value "STOPPING"
	MysqldExporterItems0StatusSTOPPING string = "STOPPING"

	// MysqldExporterItems0StatusDONE captures enum value "DONE"
	MysqldExporterItems0StatusDONE string = "DONE"
)

// prop value enum
func (o *MysqldExporterItems0) validateStatusEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, mysqldExporterItems0TypeStatusPropEnum); err != nil {
		return err
	}
	return nil
}

func (o *MysqldExporterItems0) validateStatus(formats strfmt.Registry) error {

	if swag.IsZero(o.Status) { // not required
		return nil
	}

	// value enum
	if err := o.validateStatusEnum("status", "body", *o.Status); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *MysqldExporterItems0) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *MysqldExporterItems0) UnmarshalBinary(b []byte) error {
	var res MysqldExporterItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*NodeExporterItems0 NodeExporter runs on Generic on Container Node and exposes its metrics.
swagger:model NodeExporterItems0
*/
type NodeExporterItems0 struct {

	// Unique randomly generated instance identifier.
	AgentID string `json:"agent_id,omitempty"`

	// Listen port for scraping metrics.
	ListenPort int64 `json:"listen_port,omitempty"`

	// Node identifier where this instance runs.
	NodeID string `json:"node_id,omitempty"`

	// AgentStatus represents actual Agent process status.
	// Enum: [AGENT_STATUS_INVALID STARTING RUNNING WAITING STOPPING DONE]
	Status *string `json:"status,omitempty"`
}

// Validate validates this node exporter items0
func (o *NodeExporterItems0) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateStatus(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var nodeExporterItems0TypeStatusPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["AGENT_STATUS_INVALID","STARTING","RUNNING","WAITING","STOPPING","DONE"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		nodeExporterItems0TypeStatusPropEnum = append(nodeExporterItems0TypeStatusPropEnum, v)
	}
}

const (

	// NodeExporterItems0StatusAGENTSTATUSINVALID captures enum value "AGENT_STATUS_INVALID"
	NodeExporterItems0StatusAGENTSTATUSINVALID string = "AGENT_STATUS_INVALID"

	// NodeExporterItems0StatusSTARTING captures enum value "STARTING"
	NodeExporterItems0StatusSTARTING string = "STARTING"

	// NodeExporterItems0StatusRUNNING captures enum value "RUNNING"
	NodeExporterItems0StatusRUNNING string = "RUNNING"

	// NodeExporterItems0StatusWAITING captures enum value "WAITING"
	NodeExporterItems0StatusWAITING string = "WAITING"

	// NodeExporterItems0StatusSTOPPING captures enum value "STOPPING"
	NodeExporterItems0StatusSTOPPING string = "STOPPING"

	// NodeExporterItems0StatusDONE captures enum value "DONE"
	NodeExporterItems0StatusDONE string = "DONE"
)

// prop value enum
func (o *NodeExporterItems0) validateStatusEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, nodeExporterItems0TypeStatusPropEnum); err != nil {
		return err
	}
	return nil
}

func (o *NodeExporterItems0) validateStatus(formats strfmt.Registry) error {

	if swag.IsZero(o.Status) { // not required
		return nil
	}

	// value enum
	if err := o.validateStatusEnum("status", "body", *o.Status); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *NodeExporterItems0) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *NodeExporterItems0) UnmarshalBinary(b []byte) error {
	var res NodeExporterItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*PMMAgentItems0 PMMAgent runs on Generic on Container Node.
swagger:model PMMAgentItems0
*/
type PMMAgentItems0 struct {

	// Unique randomly generated instance identifier.
	AgentID string `json:"agent_id,omitempty"`

	// True if Agent is running and connected to pmm-managed.
	Connected bool `json:"connected,omitempty"`

	// Node identifier where this instance runs.
	NodeID string `json:"node_id,omitempty"`
}

// Validate validates this PMM agent items0
func (o *PMMAgentItems0) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PMMAgentItems0) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PMMAgentItems0) UnmarshalBinary(b []byte) error {
	var res PMMAgentItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*RDSExporterItems0 RDSExporter runs on Generic or Container Node and exposes RemoteAmazonRDS Node and AmazonRDSMySQL Service metrics.
swagger:model RDSExporterItems0
*/
type RDSExporterItems0 struct {

	// Unique randomly generated instance identifier.
	AgentID string `json:"agent_id,omitempty"`

	// Listen port for scraping metrics.
	ListenPort int64 `json:"listen_port,omitempty"`

	// Node identifier where this instance runs.
	RunsOnNodeID string `json:"runs_on_node_id,omitempty"`

	// A list of Service identifiers (Node identifiers are extracted from Services).
	ServiceIds []string `json:"service_ids"`

	// AgentStatus represents actual Agent process status.
	// Enum: [AGENT_STATUS_INVALID STARTING RUNNING WAITING STOPPING DONE]
	Status *string `json:"status,omitempty"`
}

// Validate validates this RDS exporter items0
func (o *RDSExporterItems0) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateStatus(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var rdsExporterItems0TypeStatusPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["AGENT_STATUS_INVALID","STARTING","RUNNING","WAITING","STOPPING","DONE"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		rdsExporterItems0TypeStatusPropEnum = append(rdsExporterItems0TypeStatusPropEnum, v)
	}
}

const (

	// RDSExporterItems0StatusAGENTSTATUSINVALID captures enum value "AGENT_STATUS_INVALID"
	RDSExporterItems0StatusAGENTSTATUSINVALID string = "AGENT_STATUS_INVALID"

	// RDSExporterItems0StatusSTARTING captures enum value "STARTING"
	RDSExporterItems0StatusSTARTING string = "STARTING"

	// RDSExporterItems0StatusRUNNING captures enum value "RUNNING"
	RDSExporterItems0StatusRUNNING string = "RUNNING"

	// RDSExporterItems0StatusWAITING captures enum value "WAITING"
	RDSExporterItems0StatusWAITING string = "WAITING"

	// RDSExporterItems0StatusSTOPPING captures enum value "STOPPING"
	RDSExporterItems0StatusSTOPPING string = "STOPPING"

	// RDSExporterItems0StatusDONE captures enum value "DONE"
	RDSExporterItems0StatusDONE string = "DONE"
)

// prop value enum
func (o *RDSExporterItems0) validateStatusEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, rdsExporterItems0TypeStatusPropEnum); err != nil {
		return err
	}
	return nil
}

func (o *RDSExporterItems0) validateStatus(formats strfmt.Registry) error {

	if swag.IsZero(o.Status) { // not required
		return nil
	}

	// value enum
	if err := o.validateStatusEnum("status", "body", *o.Status); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *RDSExporterItems0) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *RDSExporterItems0) UnmarshalBinary(b []byte) error {
	var res RDSExporterItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
