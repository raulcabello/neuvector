package api

import (
	"github.com/neuvector/neuvector/share"
)

const (
	FedRoleNone   = ""
	FedRoleMaster = "master"
	FedRoleJoint  = "joint"
)

const (
	FedClusterStatusNone           = "active"
	FedClusterStatusCmdUnknown     = "unknown_cmd"
	FedClusterStatusCmdReceived    = "notified"
	FedClusterStatusCmdReqError    = "req_error"
	FedStatusMasterUpgradeRequired = "master_upgrade_required" // for describing master cluster only
	FedStatusJointUpgradeRequired  = "joint_upgrade_required"  // for describing joint cluster only
	FedStatusClusterUpgradeOngoing = "cluster_upgrade_ongoing" // could get this status only when rolling upgrade & polling fed rules on joint cluster are happenging
	FedStatusJointVersionTooNew    = "joint_version_too_new"   // for describing joint cluster only
	FedStatusClusterConnected      = "connected"               // for describing master cluster only
	FedStatusClusterDisconnected   = "disconnected"            // for describing master cluster only
	FedStatusClusterJoined         = "joined"                  // for describing joint cluster only. short-lived (between joining and the first polling)
	FedStatusClusterOutOfSync      = "out_of_sync"             // for describing joint cluster only
	FedStatusClusterSynced         = "synced"                  // for describing joint cluster only
	FedStatusClusterKicked         = "kicked"                  // for describing self on joint cluster only
	FedStatusClusterLeft           = "left"                    // for describing joint cluster only
	FedStatusLicenseDisallowed     = "license_disallow"        // for describing clusters in fed
	FedStatusClusterPinging        = "pinging"                 // for describing joint cluster only. short-lived (between license update and the immediate ping)
	FedStatusClusterSyncing        = "syncing"                 // for describing joint cluster only. short-lived (when joint cluster is applying fed rules)
)

// master cluster: a promoted cluster. One per-federation
// joint cluster: the other non-master clusters that join the federation
// 1. A cluster becomes master cluster after it's promoted (providing the ip/port for other clusters to join)
// 2. A cluster can join one federation at most
// 3. A master cluster with joint cluster(s) cannot join other federation
// 4. A master cluster without any joint cluster can join another federation. At the same time it becomes joint cluster of another federation.
type RESTFedMasterClusterInfo struct {
	Disabled bool                     `json:"disabled"`
	Name     string                   `json:"name"` // cluster name
	ID       string                   `json:"id"`
	Secret   string                   `json:"secret"` // used for encryoting/decrypting join_ticket issued by the master cluster. never export
	User     string                   `json:"user"`   // the user who promoets this cluster to master cluster in federation
	Status   string                   `json:"status"` // ex: FedStatusClusterSynced/FedStatusClusterOutOfSync (see above)
	RestInfo share.CLUSRestServerInfo `json:"rest_info"`
}

type RESTFedJointClusterInfo struct {
	Disabled      bool                     `json:"disabled"`
	Name          string                   `json:"name"` // cluster name
	ID            string                   `json:"id"`
	Secret        string                   `json:"secret"`
	User          string                   `json:"user"`   // the user who joins this cluster to federation
	Status        string                   `json:"status"` // ex: FedStatusClusterSynced/FedStatusClusterOutOfSync (see above)
	RestInfo      share.CLUSRestServerInfo `json:"rest_info"`
	ProxyRequired bool                     `json:"proxy_required"` // a joint cluster may be reachable without proxy even master cluster is configured to use proxy. decided when it joins fed.
}

type RESTFedMembereshipData struct { // including all clusters in the federation
	FedRole           string                     `json:"fed_role"`                 // FedRoleMaster / FedRoleJoint / FedRoleNone (see above)
	LocalRestInfo     share.CLUSRestServerInfo   `json:"local_rest_info"`          //
	MasterCluster     *RESTFedMasterClusterInfo  `json:"master_cluster,omitempty"` // master cluster
	JointClusters     []*RESTFedJointClusterInfo `json:"joint_clusters"`           // all non-master clusters in the federation
	UseProxy          string                     `json:"use_proxy"`                // http / https
	DeployRegScanData bool                       `json:"deploy_reg_scan_data"`     // whether fed registry scan data deployment is enabled
}

type RESTFedConfigData struct { // including all clusters in the federation
	PingInterval      *uint32                   `json:"ping_interval,omitempty"` // in minute
	PollInterval      *uint32                   `json:"poll_interval,omitempty"` // in minute
	Name              *string                   `json:"name,omitempty"`          // cluster name
	RestInfo          *share.CLUSRestServerInfo `json:"rest_info,omitempty"`
	UseProxy          *string                   `json:"use_proxy,omitempty"`  // http / https
	DeployRegScanData *bool                     `json:"deploy_reg_scan_data"` // whether fed registry scan data deployment is enabled
}

type RESTFedPromoteReqData struct {
	Name              string                    `json:"name,omitempty"`             // cluster name
	PingInterval      uint32                    `json:"ping_interval"`              // in minute
	PollInterval      uint32                    `json:"poll_interval"`              // in minute
	MasterRestInfo    *share.CLUSRestServerInfo `json:"master_rest_info,omitempty"` // rest info about this master cluster
	UseProxy          *string                   `json:"use_proxy,omitempty"`        // http / https
	DeployRegScanData *bool                     `json:"deploy_reg_scan_data"`       // whether fed registry scan data deployment is enabled
}

type RESTFedPromoteRespData struct {
	FedRole           string                   `json:"fed_role"`
	MasterCluster     RESTFedMasterClusterInfo `json:"master_cluster"`       // info about this master cluster
	UseProxy          string                   `json:"use_proxy,omitempty"`  // http / https
	DeployRegScanData bool                     `json:"deploy_reg_scan_data"` // whether fed registry scan data deployment is enabled
}

type RESTFedJoinToken struct { // json of the join token that contains master cluster server/port & encrypted join_ticket
	JoinToken string `json:"join_token"`
}

type RESTFedJoinReq struct { // from manager to local cluster
	Name          string                    `json:"name"`                      // cluster name
	Server        string                    `json:"server"`                    // server of master cluster
	Port          uint                      `json:"port"`                      // port of master cluster
	JoinToken     string                    `json:"join_token"`                // generated by the master cluster, i.e. RESTFedJoinToken.JoinToken
	JointRestInfo *share.CLUSRestServerInfo `json:"joint_rest_info,omitempty"` // rest info about this joint cluster
	UseProxy      *string                   `json:"use_proxy,omitempty"`
}

type RESTFedJoinReqInternal struct { // from joining cluster to master cluster for joining federation.
	User         string                  `json:"user"`                   // current operating user
	Remote       string                  `json:"remote"`                 // current operating user's remote info
	UserRoles    map[string]string       `json:"user_roles"`             // current operating user's roles
	FedKvVersion string                  `json:"fed_kv_version"`         // kv version in the code of the joining cluster
	RestVersion  string                  `json:"rest_version,omitempty"` // rest version in the code of joining cluster
	JoinTicket   string                  `json:"join_ticket"`            // generated by the master cluster, not containing master's server/port
	JointCluster RESTFedJointClusterInfo `json:"joint_cluster"`          // info about joint cluster
}

type RESTFedJoinRespInternal struct { // from master cluster to joining cluster for response of joining federation.
	PollInterval  uint32                    `json:"poll_interval"`  // in minute
	CACert        string                    `json:"ca_cert"`        // ca cert for the federated rest server in master cluster
	ClientKey     string                    `json:"client_key"`     // client key for the joint cluster
	ClientCert    string                    `json:"client_cert"`    // client cert for the joint cluster
	MasterCluster *RESTFedMasterClusterInfo `json:"master_cluster"` // info about the master cluster
}

type RESTFedLeaveReq struct { // from manager to joint cluster
	Force bool `json:"force"` // true means leave federation no matter master cluster succeeds or not
}

// for leaving federation request from joint clusters to master cluster
type RESTFedLeaveReqInternal struct { // from joint cluster to master cluster for leaving federation.
	ID          string            `json:"id"`           // id of the joint cluster to leave federation
	JointTicket string            `json:"joint_ticket"` // generated using joint cluster's secret
	User        string            `json:"user"`         // current operating user
	Remote      string            `json:"remote"`       // current operating user's remote info
	UserRoles   map[string]string `json:"user_roles"`   // current operating user's roles
}

type RESTFedRemovedReqInternal struct { // from master cluster to joint cluster for being removed from federation.
	User string `json:"user"` // current operating user
}

type RESTFedTokenResp struct {
	Token string `json:"token"` // for issued by remote joint cluster
}

// for deploying fed settings to joint clusters
type RESTDeployFedRulesReq struct {
	Force bool     `json:"force"` // true means deploying all federal rules. false means only deploying the newly changed federal rules.
	IDs   []string `json:"ids"`   // empty means deploy to all clusters
}

type RESTDeployFedRulesResp struct {
	Results map[string]int `json:"results"` // value: _fedSuccess/....
}

type RESTFedRulesSettings struct {
	AdmCtrlRulesData    *share.CLUSFedAdmCtrlRulesData   `json:"admctrl_rules_data,omitempty"`
	NetworkRulesData    *share.CLUSFedNetworkRulesData   `json:"network_rules_data,omitempty"`
	ResponseRulesData   *share.CLUSFedResponseRulesData  `json:"response_rules_data,omitempty"`
	GroupsData          *share.CLUSFedGroupsData         `json:"groups_data,omitempty"`
	FileMonitorData     *share.CLUSFedFileMonitorData    `json:"file_monitor_data,omitempty"`
	ProcessProfilesData *share.CLUSFedProcessProfileData `json:"process_profiles_data,omitempty"`
	SystemConfigData    *share.CLUSFedSystemConfigData   `json:"system_config_data,omitempty"`
}

type RESTFedImageScanResult struct {
	Revision string                          `json:"revision"` // it's md5 of json.marshal(Summary+Report)
	Summary  *share.CLUSRegistryImageSummary `json:"summary,omitempty"`
	Report   *share.CLUSScanReport           `json:"report,omitempty"`
}

type RESTFedScanData struct {
	UpdatedScanResults map[string]map[string]*RESTFedImageScanResult `json:"updated_scan_result,omitempty"` // registry name : image id : scan result; it contains only new/updated scan results
	DeletedScanResults map[string][]string                           `json:"deleted_scan_result,omitempty"` // registry name : []image id. map value being nil means the registry is deleted
}

type RESTFedInternalCommandReq struct {
	FedKvVersion string            `json:"fed_kv_version"` // kv version in the code of master cluster
	Command      string            `json:"command"`        // currently supported commands: _cmdPollFedRules / _cmdForcePullFedRules
	User         string            `json:"user"`           // current operating user
	Revisions    map[string]uint64 `json:"revisions"`      // key is fed rules type, value is the revision of current fed rules
}

type RESTFedInternalCommandResp struct {
	Result int `json:"result"` // value: _fedCmdReceived/....
}

type RESTFedPingReq struct { // from manager to joint cluster
	Token        string `json:"token"`
	FedKvVersion string `json:"fed_kv_version"` // kv version in the code of the master cluster
}

type RESTFedPingResp struct { // from manager to joint cluster
	Result int `json:"result"` // value: _fedSuccess/....
}

// for polling fed rules/settings from joint clusters to master cluster
type RESTPollFedRulesReq struct {
	ID                  string            `json:"id"`                          // id of joint cluster
	Name                string            `json:"name"`                        // name of joint cluster
	JointTicket         string            `json:"joint_ticket"`                // generated using joint cluster's secret
	FedKvVersion        string            `json:"fed_kv_version"`              // kv version in the code of joint cluster
	RestVersion         string            `json:"rest_version,omitempty"`      // rest version in the code of joint cluster
	Revisions           map[string]uint64 `json:"revisions"`                   // key is fed rules type, value is the revision
	RegistryRevision    uint64            `json:"registry_revision"`           // fed registry revision
	ScannedRegImagesRev uint64            `json:"scanned_reg_images_revision"` // scanned images' revision in fed registries on master cluster
}

type RESTPollFedRulesResp struct {
	Result              int               `json:"result"`                      // value: _fedSuccess/....
	PollInterval        uint32            `json:"poll_interval"`               // in minute
	Settings            []byte            `json:"settings,omitempty"`          // marshall of RESTFedRulesSettings, which contains only modified settings (for ~5.0.x)
	Revisions           map[string]uint64 `json:"revisions"`                   // key is fed rules type, value is the revision. It contains only revisions of modified settings
	RegistryRevision    uint64            `json:"registry_revision"`           // fed registry revision
	ScannedRegImagesRev uint64            `json:"scanned_reg_images_revision"` // scanned images' revision in fed registries on master cluster
	DeployRegScanData   bool              `json:"deploy_reg_scan_data"`        // for informing whether master cluster accepts registry scan data polling
}

type RESTPollFedScanDataReq struct {
	ID                  string                       `json:"id"`                          // id of joint cluster
	Name                string                       `json:"name"`                        // name of joint cluster
	JointTicket         string                       `json:"joint_ticket"`                // generated using joint cluster's secret
	FedKvVersion        string                       `json:"fed_kv_version"`              // kv version in the code of joint cluster
	RestVersion         string                       `json:"rest_version,omitempty"`      // rest version in the code of joint cluster
	RegistryRevision    uint64                       `json:"registry_revision"`           // fed registry revision
	ScannedRegImagesRev uint64                       `json:"scanned_reg_images_revision"` // scanned images' revision in fed registries on master cluster that the managed cluster remembers
	ScanReportRevisions map[string]map[string]string `json:"scan_report_revisions"`       // registry name : image id : scan report md5
}

type RESTPollFedScanDataResp struct {
	Result              int                          `json:"result"`                      // value: _fedSuccess/....
	PollInterval        uint32                       `json:"poll_interval"`               // in minute
	RegistryRevision    uint64                       `json:"registry_revision"`           // fed registry revision
	RegistryData        *share.CLUSFedRegistriesData `json:"registry_data,omitempty"`     // all fed registry settings if there is any change since last polling
	ScannedRegImagesRev uint64                       `json:"scanned_reg_images_revision"` // scanned images' revision in fed registries on master cluster
	FedScanData         RESTFedScanData              `json:"fed_scan_data"`               // (partial) updated/deleted scan result since last polling
	HasPartialScanData  bool                         `json:"has_partial_scan_data"`       // true when master cluster returns partial new scan data. managed clusters should keep polling based on bandwidth consideration.
	ThrottleTime        uint32                       `json:"throttle_time"`               // in ms
}

type RESTFedView struct {
	Compatible bool `json:"compatible"`
}
