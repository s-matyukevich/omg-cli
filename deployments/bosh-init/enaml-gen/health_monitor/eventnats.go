package health_monitor 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type EventNats struct {

	/*Address - Descr: Address of the event NATS message bus to connect to Default: <nil>
*/
	Address interface{} `yaml:"address,omitempty"`

	/*Password - Descr: Password for event NATS message bus connection Default: <nil>
*/
	Password interface{} `yaml:"password,omitempty"`

	/*User - Descr: User for the event NATS message bus connection Default: <nil>
*/
	User interface{} `yaml:"user,omitempty"`

	/*Port - Descr: Port of the event NATS message bus port to connect to Default: <nil>
*/
	Port interface{} `yaml:"port,omitempty"`

}