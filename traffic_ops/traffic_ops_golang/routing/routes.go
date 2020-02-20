package routing

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"encoding/json"
	"errors"
	"net/http"
	"runtime"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/about"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/apicapability"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/apiriak"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/apitenant"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/asn"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/cachegroup"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/cachegroupparameter"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/cachesstats"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/capabilities"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/cdn"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/cdnfederation"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/coordinate"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/crconfig"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbdump"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/deliveryservice"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/deliveryservice/consistenthash"
	dsrequest "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/deliveryservice/request"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/deliveryservice/request/comment"
	dsserver "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/deliveryservice/servers"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/deliveryservicerequests"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/deliveryservicesregexes"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/division"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/federation_resolvers"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/federations"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/invalidationjobs"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/iso"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/login"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/logs"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/origin"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/parameter"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/physlocation"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/ping"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/profile"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/profileparameter"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/region"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/role"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/server"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/servercapability"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/servercheck"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/staticdnsentry"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/status"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/steering"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/steeringtargets"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/systeminfo"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/toextension"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/trafficstats"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/types"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/urisigning"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/user"

	"github.com/jmoiron/sqlx"
)

// Authenticated ...
const Authenticated = true

// NoAuth ...
const NoAuth = false

func handlerToFunc(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}
}

// Routes returns the API routes, raw non-API root level routes, and a catchall route for when no route matches.
func Routes(d ServerData) ([]Route, []RawRoute, http.HandlerFunc, error) {
	routes := []Route{

		// API Capability
		{api.Version{2, 0}, http.MethodGet, `api_capabilities/?$`, apicapability.GetAPICapabilitiesHandler, auth.PrivLevelReadOnly, Authenticated, nil},

		//ASNs
		{api.Version{2, 0}, http.MethodPut, `asns/?$`, api.UpdateHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `asns/?$`, api.DeleteHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `asns/?$`, api.ReadHandler(&asn.TOASNV11{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `asns/{id}$`, api.ReadHandler(&asn.TOASNV11{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPut, `asns/{id}$`, api.UpdateHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `asns/?$`, api.CreateHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `asns/{id}$`, api.DeleteHandler(&asn.TOASNV11{}), auth.PrivLevelOperations, Authenticated, nil},

		// Traffic Stats access
		{api.Version{2, 0}, http.MethodGet, `deliveryservice_stats`, trafficstats.GetDSStats, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `cache_stats`, trafficstats.GetCacheStats, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `current_stats/?$`, trafficstats.GetCurrentStats, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `caches/stats/?$`, cachesstats.Get, auth.PrivLevelReadOnly, Authenticated, nil},

		//Cache Groups manipulations
		{api.Version{2, 0}, http.MethodGet, `cachegroups/trimmed/?$`, cachegroup.GetTrimmed, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `cachegroups/?$`, api.ReadHandler(&cachegroup.TOCacheGroup{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `cachegroups/{id}$`, api.ReadHandler(&cachegroup.TOCacheGroup{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPut, `cachegroups/{id}$`, api.UpdateHandler(&cachegroup.TOCacheGroup{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `cachegroups/?$`, api.CreateHandler(&cachegroup.TOCacheGroup{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `cachegroups/{id}$`, api.DeleteHandler(&cachegroup.TOCacheGroup{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `cachegroups/{id}/queue_update$`, cachegroup.QueueUpdates, auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `cachegroups/{id}/deliveryservices/?$`, cachegroup.DSPostHandler, auth.PrivLevelOperations, Authenticated, nil},

		// Cache-Group-Parameter associations
		{api.Version{2, 0}, http.MethodGet, `cachegroupparameters/?$`, cachegroupparameter.ReadAllCacheGroupParameters, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `cachegroupparameters/?$`, cachegroupparameter.AddCacheGroupParameters, auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `cachegroups/{id}/parameters/?$`, api.ReadHandler(&cachegroupparameter.TOCacheGroupParameter{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `cachegroups/{id}/unassigned_parameters/?$`, api.ReadHandler(&cachegroupparameter.TOCacheGroupUnassignedParameter{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `cachegroupparameters/{cachegroupID}/{parameterId}$`, api.DeleteHandler(&cachegroupparameter.TOCacheGroupParameter{}), auth.PrivLevelOperations, Authenticated, nil},

		// User Capabilities
		{api.Version{2, 0}, http.MethodGet, `capabilities/?$`, capabilities.Read, auth.PrivLevelReadOnly, Authenticated, nil},

		//CDN
		{api.Version{2, 0}, http.MethodGet, `cdns/metric_types`, notImplementedHandler, 0, NoAuth, nil},
		{api.Version{2, 0}, http.MethodGet, `cdns/capacity$`, cdn.GetCapacity, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `cdns/{name}/health/?$`, cdn.GetNameHealth, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `cdns/health/?$`, cdn.GetHealth, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `cdns/domains/?$`, cdn.DomainsHandler, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `cdns/routing$`, cdn.GetRouting, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `cdns/name/{name}/?$`, api.ReadHandler(&cdn.TOCDN{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `cdns/name/{name}$`, cdn.DeleteName, auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `cdns/{id}/queue_update$`, cdn.Queue, auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `cdns/?$`, api.ReadHandler(&cdn.TOCDN{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `cdns/{id}$`, api.ReadHandler(&cdn.TOCDN{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPut, `cdns/{id}$`, api.UpdateHandler(&cdn.TOCDN{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `cdns/?$`, api.CreateHandler(&cdn.TOCDN{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `cdns/{id}$`, api.DeleteHandler(&cdn.TOCDN{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `cdns/{cdn}/configs/monitoring?$`, crconfig.SnapshotGetMonitoringHandler, auth.PrivLevelReadOnly, Authenticated, nil},

		// CDN SSL/DNSSEC Keys
		{api.Version{2, 0}, http.MethodGet, `cdns/name/{name}/sslkeys/?$`, cdn.GetSSLKeys, auth.PrivLevelAdmin, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `cdns/dnsseckeys/generate?$`, cdn.CreateDNSSECKeys, auth.PrivLevelAdmin, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `cdns/name/{name}/dnsseckeys/delete/?$`, cdn.DeleteDNSSECKeys, auth.PrivLevelAdmin, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `cdns/name/{name}/dnsseckeys/?$`, cdn.GetDNSSECKeys, auth.PrivLevelAdmin, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `cdns/dnsseckeys/refresh/?$`, cdn.RefreshDNSSECKeys, auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `cdns/{name}/dnsseckeys/ksk/generate$`, cdn.GenerateKSK, auth.PrivLevelAdmin, Authenticated, nil},

		//Database dumps
		{api.Version{2, 0}, http.MethodGet, `dbdump/?`, dbdump.DBDump, auth.PrivLevelAdmin, Authenticated, nil},

		//Division
		{api.Version{2, 0}, http.MethodGet, `divisions/??$`, api.ReadHandler(&division.TODivision{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `divisions/{id}$`, api.ReadHandler(&division.TODivision{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPut, `divisions/{id}$`, api.UpdateHandler(&division.TODivision{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `divisions/?$`, api.CreateHandler(&division.TODivision{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `divisions/{id}$`, api.DeleteHandler(&division.TODivision{}), auth.PrivLevelOperations, Authenticated, nil},

		// Changelogs
		{api.Version{2, 0}, http.MethodGet, `logs/?$`, logs.Get, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `logs/{days}/days/?$`, logs.Get, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `logs/newcount/?$`, logs.GetNewCount, auth.PrivLevelReadOnly, Authenticated, nil},

		//Content invalidation jobs
		{api.Version{2, 0}, http.MethodGet, `jobs/?$`, api.ReadHandler(&invalidationjobs.InvalidationJob{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `jobs/?$`, invalidationjobs.Delete, auth.PrivLevelPortal, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPut, `jobs/?$`, invalidationjobs.Update, auth.PrivLevelPortal, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `jobs/?`, invalidationjobs.Create, auth.PrivLevelPortal, Authenticated, nil},

		//Login
		{api.Version{2, 0}, http.MethodGet, `users/{id}/deliveryservices/?$`, user.GetDSes, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `user/{id}/deliveryservices/available/?$`, user.GetAvailableDSes, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `user/login/?$`, login.LoginHandler(d.DB, d.Config), 0, NoAuth, nil},
		{api.Version{2, 0}, http.MethodPost, `user/logout/?$`, login.LogoutHandler(d.Config.Secrets[0]), 0, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `user/login/oauth/?$`, login.OauthLoginHandler(d.DB, d.Config), 0, NoAuth, nil},
		{api.Version{2, 0}, http.MethodPost, `user/login/token/?$`, login.TokenLoginHandler(d.DB, d.Config), 0, NoAuth, nil},
		{api.Version{2, 0}, http.MethodPost, `user/reset_password/?$`, login.ResetPassword(d.DB, d.Config), 0, NoAuth, nil},
		{api.Version{2, 0}, http.MethodPost, `users/register/?$`, login.RegisterUser, auth.PrivLevelOperations, Authenticated, nil},

		//ISO
		{api.Version{2, 0}, http.MethodGet, `osversions/?$`, iso.GetOSVersions, auth.PrivLevelReadOnly, Authenticated, nil},

		//User
		{api.Version{2, 0}, http.MethodGet, `users/?$`, api.ReadHandler(&user.TOUser{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `users/{id}$`, api.ReadHandler(&user.TOUser{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPut, `users/{id}$`, api.UpdateHandler(&user.TOUser{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `users/?$`, api.CreateHandler(&user.TOUser{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `user/current/?$`, user.Current, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPut, `user/current/?$`, user.ReplaceCurrent, auth.PrivLevelReadOnly, Authenticated, nil},


		//Physical Locations
		{api.Version{2, 0}, http.MethodGet, `phys_locations/?$`, api.ReadHandler(&physlocation.TOPhysLocation{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `phys_locations/trimmed/?$`, physlocation.GetTrimmed, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `phys_locations/{id}$`, api.ReadHandler(&physlocation.TOPhysLocation{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPut, `phys_locations/{id}$`, api.UpdateHandler(&physlocation.TOPhysLocation{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `phys_locations/?$`, api.CreateHandler(&physlocation.TOPhysLocation{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `phys_locations/{id}$`, api.DeleteHandler(&physlocation.TOPhysLocation{}), auth.PrivLevelOperations, Authenticated, nil},

		//Ping
		{api.Version{2, 0}, http.MethodGet, `ping$`, ping.PingHandler(), 0, NoAuth, nil},
		{api.Version{2, 0}, http.MethodGet, `riak/ping/?$`, ping.Riak, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `keys/ping/?$`, ping.Keys, auth.PrivLevelReadOnly, Authenticated, nil},

		//Parameter
		{api.Version{2, 0}, http.MethodGet, `parameters/?$`, api.ReadHandler(&parameter.TOParameter{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `parameters/{id}$`, api.ReadHandler(&parameter.TOParameter{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPut, `parameters/{id}$`, api.UpdateHandler(&parameter.TOParameter{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `parameters/?$`, api.CreateHandler(&parameter.TOParameter{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `parameters/{id}$`, api.DeleteHandler(&parameter.TOParameter{}), auth.PrivLevelOperations, Authenticated, nil},

		// Profile-Parameter associations
		{api.Version{2, 0}, http.MethodGet, `profiles/{id}/parameters/?$`, profileparameter.GetProfileID, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `profiles/{id}/unassigned_parameters/?$`, profileparameter.GetUnassigned, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `profiles/name/{name}/parameters/?$`, profileparameter.GetProfileName, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `parameters/profile/{name}/?$`, profileparameter.GetProfileName, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `profiles/name/{name}/parameters/?$`, profileparameter.PostProfileParamsByName, auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `profiles/{id}/parameters/?$`, profileparameter.PostProfileParamsByID, auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `profileparameters/?$`, api.ReadHandler(&profileparameter.TOProfileParameter{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `profileparameters/?$`, api.CreateHandler(&profileparameter.TOProfileParameter{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `profileparameter/?$`, profileparameter.PostProfileParam, auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `parameterprofile/?$`, profileparameter.PostParamProfile, auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `profileparameters/{profileId}/{parameterId}$`, api.DeleteHandler(&profileparameter.TOProfileParameter{}), auth.PrivLevelOperations, Authenticated, nil},

		//Profile
		{api.Version{2, 0}, http.MethodGet, `profiles/?$`, api.ReadHandler(&profile.TOProfile{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `profiles/trimmed/?$`, profile.Trimmed, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `profiles/{id}$`, api.ReadHandler(&profile.TOProfile{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPut, `profiles/{id}$`, api.UpdateHandler(&profile.TOProfile{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `profiles/?$`, api.CreateHandler(&profile.TOProfile{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `profiles/{id}$`, api.DeleteHandler(&profile.TOProfile{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `profiles/{id}/export/?$`, profile.ExportProfileHandler, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `profiles/import/?$`, profile.ImportProfileHandler, auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `profiles/name/{new_profile}/copy/{existing_profile}`, profile.CopyProfileHandler, auth.PrivLevelOperations, Authenticated, nil},

		//Region
		{api.Version{2, 0}, http.MethodGet, `regions/?$`, api.ReadHandler(&region.TORegion{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `regions/{id}$`, api.ReadHandler(&region.TORegion{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `regions/name/{name}/?$`, region.GetName, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPut, `regions/{id}$`, api.UpdateHandler(&region.TORegion{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `regions/?$`, api.CreateHandler(&region.TORegion{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `regions/?$`, api.DeleteHandler(&region.TORegion{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `regions/{id}$`, api.DeleteHandler(&region.TORegion{}), auth.PrivLevelOperations, Authenticated, nil},

		// DeliveryServices
		{api.Version{2, 0}, http.MethodGet, `deliveryservices/?$`, api.ReadHandler(&deliveryservice.TODeliveryService{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `deliveryservices/{id}/?$`, api.ReadHandler(&deliveryservice.TODeliveryService{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `deliveryservices/?$`, deliveryservice.CreateV15, auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPut, `deliveryservices/{id}/?$`, deliveryservice.UpdateV15, auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPut, `deliveryservices/{id}/safe/?$`, deliveryservice.UpdateSafe, auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `deliveryservices/{id}/?$`, api.DeleteHandler(&deliveryservice.TODeliveryService{}), auth.PrivLevelOperations, Authenticated, nil},

		// Delivery Service Requests
		{api.Version{2, 0}, http.MethodPost, `deliveryservices/request`, deliveryservicerequests.Request, auth.PrivLevelPortal, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `deliveryservice_requests/?$`, api.ReadHandler(&dsrequest.TODeliveryServiceRequest{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPut, `deliveryservice_requests/?$`, api.UpdateHandler(&dsrequest.TODeliveryServiceRequest{}), auth.PrivLevelPortal, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `deliveryservice_requests/?$`, api.CreateHandler(&dsrequest.TODeliveryServiceRequest{}), auth.PrivLevelPortal, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `deliveryservice_requests/?$`, api.DeleteHandler(&dsrequest.TODeliveryServiceRequest{}), auth.PrivLevelPortal, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPut, `deliveryservice_requests/{id}/assign$`, api.UpdateHandler(dsrequest.GetAssignmentSingleton()), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPut, `deliveryservice_requests/{id}/status$`, api.UpdateHandler(dsrequest.GetStatusSingleton()), auth.PrivLevelPortal, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `deliveryservice_request_comments/?$`, api.ReadHandler(&comment.TODeliveryServiceRequestComment{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPut, `deliveryservice_request_comments/?$`, api.UpdateHandler(&comment.TODeliveryServiceRequestComment{}), auth.PrivLevelPortal, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `deliveryservice_request_comments/?$`, api.CreateHandler(&comment.TODeliveryServiceRequestComment{}), auth.PrivLevelPortal, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `deliveryservice_request_comments/?$`, api.DeleteHandler(&comment.TODeliveryServiceRequestComment{}), auth.PrivLevelPortal, Authenticated, nil},

		//Delivery Service SSL/URI-Signing/URL-Signature/DNSSEC keys
		{api.Version{2, 0}, http.MethodGet, `deliveryservices/{xmlID}/urisignkeys$`, urisigning.GetURIsignkeysHandler, auth.PrivLevelAdmin, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `deliveryservices/{xmlID}/urisignkeys$`, urisigning.SaveDeliveryServiceURIKeysHandler, auth.PrivLevelAdmin, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPut, `deliveryservices/{xmlID}/urisignkeys$`, urisigning.SaveDeliveryServiceURIKeysHandler, auth.PrivLevelAdmin, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `deliveryservices/{xmlID}/urisignkeys$`, urisigning.RemoveDeliveryServiceURIKeysHandler, auth.PrivLevelAdmin, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `deliveryservices/xmlId/{xmlid}/sslkeys$`, deliveryservice.GetSSLKeysByXMLIDV15, auth.PrivLevelAdmin, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `deliveryservices/hostname/{hostname}/sslkeys$`, deliveryservice.GetSSLKeysByHostNameV15, auth.PrivLevelAdmin, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `deliveryservices/sslkeys/add$`, deliveryservice.AddSSLKeys, auth.PrivLevelAdmin, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `deliveryservices/xmlId/{xmlid}/sslkeys/delete$`, deliveryservice.DeleteSSLKeys, auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `deliveryservices/sslkeys/generate/?$`, deliveryservice.GenerateSSLKeys, auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `deliveryservices/xmlId/{name}/urlkeys/copyFromXmlId/{copy-name}/?$`, deliveryservice.CopyURLKeys, auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `deliveryservices/xmlId/{name}/urlkeys/generate/?$`, deliveryservice.GenerateURLKeys, auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `deliveryservices/xmlId/{name}/urlkeys/?$`, deliveryservice.GetURLKeysByName, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `deliveryservices/{id}/urlkeys/?$`, deliveryservice.GetURLKeysByID, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `riak/bucket/{bucket}/key/{key}/values/?$`, apiriak.GetBucketKey, auth.PrivLevelAdmin, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `deliveryservices/sslkeys/generate/letsencrypt/?$`, deliveryservice.GenerateLetsEncryptCertificates, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `letsencrypt/dnsrecords/?$`, deliveryservice.GetDnsChallengeRecords, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `letsencrypt/autorenew/?$`, deliveryservice.RenewCertificates, auth.PrivLevelOperations, Authenticated, nil},

		//Delivery-Service-Capabilities associations
		{api.Version{2, 0}, http.MethodGet, `deliveryservices_required_capabilities/?$`, api.ReadHandler(&deliveryservice.RequiredCapability{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `deliveryservices_required_capabilities/?$`, api.CreateHandler(&deliveryservice.RequiredCapability{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `deliveryservices_required_capabilities/?$`, api.DeleteHandler(&deliveryservice.RequiredCapability{}), auth.PrivLevelOperations, Authenticated, nil},

		// Delivery Service routing and statistics
		{api.Version{2, 0}, http.MethodGet, `deliveryservice_matches/?$`, deliveryservice.GetMatches, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `deliveryservices_regexes/?$`, deliveryservicesregexes.Get, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `deliveryservices/{dsid}/regexes/?$`, deliveryservicesregexes.DSGet, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `deliveryservices/{dsid}/regexes/{regexid}?$`, deliveryservicesregexes.DSGetID, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `deliveryservices/{dsid}/regexes/?$`, deliveryservicesregexes.Post, auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPut, `deliveryservices/{dsid}/regexes/{regexid}?$`, deliveryservicesregexes.Put, auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `deliveryservices/{dsid}/regexes/{regexid}?$`, deliveryservicesregexes.Delete, auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `deliveryservices/{id}/capacity/?$`, deliveryservice.GetCapacity, auth.PrivLevelReadOnly, Authenticated, nil},

		// Delivery-Service-Server associations
		{api.Version{2, 0}, http.MethodDelete, `deliveryservice_server/{dsid}/{serverid}`, dsserver.Delete, auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `deliveryserviceserver/?$`, dsserver.ReadDSSHandlerV14, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `deliveryserviceserver$`, dsserver.GetReplaceHandler, auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `deliveryservices/{xml_id}/servers$`, dsserver.GetCreateHandler, auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `servers/{id}/deliveryservices$`, api.ReadHandler(&dsserver.TODSSDeliveryService{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `deliveryservices/{id}/servers$`, dsserver.GetReadAssigned, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `deliveryservices/{id}/unassigned_servers$`, dsserver.GetReadUnassigned, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `deliveryservices/{id}/servers/eligible/?$`, deliveryservice.GetServersEligible, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `deliveryservices/{id}/health/?$`, deliveryservice.GetHealth, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `servers/{id}/deliveryservices$`, server.AssignDeliveryServicesToServerHandler, auth.PrivLevelOperations, Authenticated, nil},

		//Servers
		{api.Version{2, 0}, http.MethodGet, `servers/status$`, server.GetServersStatusCountsHandler, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `servers/details/?$`, server.GetDetailParamHandler, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `servers/hostname/{hostName}/details/?$`, server.GetDetailHandler, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPut, `servers/{id}/status$`, server.UpdateStatusHandler, auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `servers/{id}/queue_update$`, server.QueueUpdateHandler, auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `servers/?$`, api.ReadHandler(&server.TOServer{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `servers/{id}$`, api.ReadHandler(&server.TOServer{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPut, `servers/{id}$`, api.UpdateHandler(&server.TOServer{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `servers/?$`, api.CreateHandler(&server.TOServer{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `servers/{id}$`, api.DeleteHandler(&server.TOServer{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `servers/{host_name}/update_status$`, server.GetServerUpdateStatusHandler, auth.PrivLevelReadOnly, Authenticated, nil},

		//Server-Server-Capabilities associations
		{api.Version{2, 0}, http.MethodGet, `server_server_capabilities/?$`, api.ReadHandler(&server.TOServerServerCapability{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `server_server_capabilities/?$`, api.CreateHandler(&server.TOServerServerCapability{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `server_server_capabilities/?$`, api.DeleteHandler(&server.TOServerServerCapability{}), auth.PrivLevelOperations, Authenticated, nil},

		//Server Capability
		{api.Version{2, 0}, http.MethodGet, `server_capabilities$`, api.ReadHandler(&servercapability.TOServerCapability{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `server_capabilities$`, api.CreateHandler(&servercapability.TOServerCapability{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `server_capabilities$`, api.DeleteHandler(&servercapability.TOServerCapability{}), auth.PrivLevelOperations, Authenticated, nil},

		//Serverchecks
		{api.Version{2, 0}, http.MethodGet, `servers/checks/?$`, servercheck.ReadServersChecks, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `servercheck/?$`, servercheck.CreateUpdateServercheck, auth.PrivLevelInvalid, Authenticated, nil},

		//Status
		{api.Version{2, 0}, http.MethodGet, `statuses/?$`, api.ReadHandler(&status.TOStatus{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `statuses/{id}$`, api.ReadHandler(&status.TOStatus{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPut, `statuses/{id}$`, api.UpdateHandler(&status.TOStatus{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `statuses/?$`, api.CreateHandler(&status.TOStatus{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `statuses/{id}$`, api.DeleteHandler(&status.TOStatus{}), auth.PrivLevelOperations, Authenticated, nil},

		//System Info
		{api.Version{2, 0}, http.MethodGet, `system/info/?$`, systeminfo.Get, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `about/?$`, about.Handler(), auth.PrivLevelReadOnly, Authenticated, nil},

		//Type
		{api.Version{2, 0}, http.MethodGet, `types/?$`, api.ReadHandler(&types.TOType{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `types/{id}$`, api.ReadHandler(&types.TOType{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPut, `types/{id}$`, api.UpdateHandler(&types.TOType{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `types/?$`, api.CreateHandler(&types.TOType{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `types/{id}$`, api.DeleteHandler(&types.TOType{}), auth.PrivLevelOperations, Authenticated, nil},

		//Coordinates
		{api.Version{2, 0}, http.MethodGet, `coordinates/?$`, api.ReadHandler(&coordinate.TOCoordinate{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPut, `coordinates/?$`, api.UpdateHandler(&coordinate.TOCoordinate{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `coordinates/?$`, api.CreateHandler(&coordinate.TOCoordinate{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `coordinates/?$`, api.DeleteHandler(&coordinate.TOCoordinate{}), auth.PrivLevelOperations, Authenticated, nil},

		// Federations
		{api.Version{2, 0}, http.MethodGet, `cdns/{name}/federations/?$`, api.ReadHandler(&cdnfederation.TOCDNFederation{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `cdns/{name}/federations/{id}$`, api.ReadHandler(&cdnfederation.TOCDNFederation{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `cdns/{name}/federations/?$`, api.CreateHandler(&cdnfederation.TOCDNFederation{}), auth.PrivLevelAdmin, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPut, `cdns/{name}/federations/{id}$`, api.UpdateHandler(&cdnfederation.TOCDNFederation{}), auth.PrivLevelAdmin, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `cdns/{name}/federations/{id}$`, api.DeleteHandler(&cdnfederation.TOCDNFederation{}), auth.PrivLevelAdmin, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `federations/all/?$`, federations.GetAll, auth.PrivLevelAdmin, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `federations/?$`, federations.Get, auth.PrivLevelFederation, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `federations/?$`, federations.AddFederationResolverMappingsForCurrentUser, auth.PrivLevelFederation, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `federations/?$`, federations.RemoveFederationResolverMappingsForCurrentUser, auth.PrivLevelFederation, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPut, `federations/?$`, federations.ReplaceFederationResolverMappingsForCurrentUser, auth.PrivLevelFederation, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `federations/{id}/deliveryservices/?$`, federations.PostDSes, auth.PrivLevelAdmin, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `federations/{id}/deliveryservices/?$`, api.ReadHandler(&federations.TOFedDSes{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `federations/{id}/deliveryservices/{dsID}/?$`, api.DeleteHandler(&federations.TOFedDSes{}), auth.PrivLevelAdmin, Authenticated, nil},

		// Federation Resolvers
		{api.Version{2, 0}, http.MethodPost, `federation_resolvers/?$`, federation_resolvers.Create, auth.PrivLevelAdmin, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `federation_resolvers/?$`, federation_resolvers.Read, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `federations/{id}/federation_resolvers/?$`, federations.AssignFederationResolversToFederationHandler, auth.PrivLevelAdmin, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `federations/{id}/federation_resolvers/?$`, federations.GetFederationFederationResolversHandler, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `federation_resolvers/?$`, federation_resolvers.Delete, auth.PrivLevelAdmin, Authenticated, nil},

		// Federation-User associations
		{api.Version{2, 0}, http.MethodPost, `federations/{id}/users/?$`, federations.PostUsers, auth.PrivLevelAdmin, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `federations/{id}/users/?$`, api.ReadHandler(&federations.TOUsers{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `federations/{id}/users/{userID}/?$`, api.DeleteHandler(&federations.TOUsers{}), auth.PrivLevelAdmin, Authenticated, nil},

		//Origins
		{api.Version{2, 0}, http.MethodGet, `origins/?$`, api.ReadHandler(&origin.TOOrigin{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPut, `origins/?$`, api.UpdateHandler(&origin.TOOrigin{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `origins/?$`, api.CreateHandler(&origin.TOOrigin{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `origins/?$`, api.DeleteHandler(&origin.TOOrigin{}), auth.PrivLevelOperations, Authenticated, nil},

		//Roles
		{api.Version{2, 0}, http.MethodGet, `roles/?$`, api.ReadHandler(&role.TORole{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPut, `roles/?$`, api.UpdateHandler(&role.TORole{}), auth.PrivLevelAdmin, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `roles/?$`, api.CreateHandler(&role.TORole{}), auth.PrivLevelAdmin, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `roles/?$`, api.DeleteHandler(&role.TORole{}), auth.PrivLevelAdmin, Authenticated, nil},

		//StaticDNSEntries
		{api.Version{2, 0}, http.MethodGet, `staticdnsentries/?$`, api.ReadHandler(&staticdnsentry.TOStaticDNSEntry{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPut, `staticdnsentries/?$`, api.UpdateHandler(&staticdnsentry.TOStaticDNSEntry{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `staticdnsentries/?$`, api.CreateHandler(&staticdnsentry.TOStaticDNSEntry{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `staticdnsentries/?$`, api.DeleteHandler(&staticdnsentry.TOStaticDNSEntry{}), auth.PrivLevelOperations, Authenticated, nil},

		//Tenants
		{api.Version{2, 0}, http.MethodGet, `tenants/?$`, api.ReadHandler(&apitenant.TOTenant{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `tenants/{id}$`, api.ReadHandler(&apitenant.TOTenant{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPut, `tenants/{id}$`, api.UpdateHandler(&apitenant.TOTenant{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `tenants/?$`, api.CreateHandler(&apitenant.TOTenant{}), auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `tenants/{id}$`, api.DeleteHandler(&apitenant.TOTenant{}), auth.PrivLevelOperations, Authenticated, nil},

		//Snapshots
		{api.Version{2, 0}, http.MethodGet, `cdns/{cdn}/snapshot/?$`, crconfig.SnapshotGetHandler, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `cdns/{cdn}/snapshot/new/?$`, crconfig.Handler, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPut, `cdns/{id}/snapshot/?$`, crconfig.SnapshotHandler, auth.PrivLevelOperations, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPut, `snapshot/{cdn}/?$`, crconfig.SnapshotHandler, auth.PrivLevelOperations, Authenticated, nil},


		// Steering
		{api.Version{2, 0}, http.MethodGet, `steering/{deliveryservice}/targets/?$`, api.ReadHandler(&steeringtargets.TOSteeringTargetV11{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `steering/{deliveryservice}/targets/{target}$`, api.ReadHandler(&steeringtargets.TOSteeringTargetV11{}), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `steering/{deliveryservice}/targets/?$`, api.CreateHandler(&steeringtargets.TOSteeringTargetV11{}), auth.PrivLevelSteering, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPut, `steering/{deliveryservice}/targets/{target}/?$`, api.UpdateHandler(&steeringtargets.TOSteeringTargetV11{}), auth.PrivLevelSteering, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `steering/{deliveryservice}/targets/{target}/?$`, api.DeleteHandler(&steeringtargets.TOSteeringTargetV11{}), auth.PrivLevelSteering, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `steering/?$`, steering.Get, auth.PrivLevelSteering, Authenticated, nil},

		// Stats Summary
		{api.Version{2, 0}, http.MethodGet, `stats_summary/?$`, trafficstats.GetStatsSummary, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodPost, `stats_summary/?$`, trafficstats.CreateStatsSummary, auth.PrivLevelReadOnly, Authenticated, nil},

		// TO Extensions
		{api.Version{2, 0}, http.MethodPost, `to_extensions$`, toextension.CreateTOExtension, auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodGet, `to_extensions$`, toextension.GetTOExtensionsHandler(d.Plugins), auth.PrivLevelReadOnly, Authenticated, nil},
		{api.Version{2, 0}, http.MethodDelete, `to_extensions/{id}$`, toextension.Delete, auth.PrivLevelReadOnly, Authenticated, nil},

		//Pattern based consistent hashing endpoint
		{api.Version{2, 0}, http.MethodPost, `consistenthash/?$`, consistenthash.Post, auth.PrivLevelReadOnly, Authenticated, nil},
	}

	// rawRoutes are served at the root path. These should be almost exclusively old Perl pre-API routes, which have yet to be converted in all clients. New routes should be in the versioned API path.
	rawRoutes := []RawRoute{
		// DEPRECATED - use PUT /api/1.2/snapshot/{cdn}
		{http.MethodGet, `tools/write_crconfig/{cdn}/?$`, crconfig.SnapshotOldGUIHandler, auth.PrivLevelOperations, Authenticated, nil},
		// DEPRECATED - use GET /api/1.2/cdns/{cdn}/snapshot
		{http.MethodGet, `CRConfig-Snapshots/{cdn}/CRConfig.json?$`, crconfig.SnapshotOldGetHandler, auth.PrivLevelReadOnly, Authenticated, nil},
	}

	return routes, rawRoutes, rootHandler, nil
}

func MemoryStatsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErrs := tc.GetHandleErrorsFunc(w, r)
		stats := runtime.MemStats{}
		runtime.ReadMemStats(&stats)

		bytes, err := json.Marshal(stats)
		if err != nil {
			log.Errorln("unable to marshal stats: " + err.Error())
			handleErrs(http.StatusInternalServerError, errors.New("marshalling error"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(bytes)
	}
}

func DBStatsHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErrs := tc.GetHandleErrorsFunc(w, r)
		stats := db.DB.Stats()

		bytes, err := json.Marshal(stats)
		if err != nil {
			log.Errorln("unable to marshal stats: " + err.Error())
			handleErrs(http.StatusInternalServerError, errors.New("marshalling error"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(bytes)
	}
}

const NOT_FOUND_RESPONSE = `{"alerts":[{"level":"error","text":"The requested API path was not found."}]}`

// rootHandler is the handler for all routes not found in the route definitions.
func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(rfc.ContentType, rfc.MIME_JSON.String())
	w.WriteHeader(http.StatusNotFound)
	w.Write(append([]byte(NOT_FOUND_RESPONSE), '\n'))
}

// notImplementedHandler returns a 501 Not Implemented to the client. This should be used very rarely, and primarily for old API Perl routes which were broken long ago, which we don't have the resources to rewrite in Go for the time being.
func notImplementedHandler(w http.ResponseWriter, r *http.Request) {
	code := http.StatusNotImplemented
	w.WriteHeader(code)
	w.Write([]byte(http.StatusText(code)))
}

//CreateThrottledHandler takes a handler, and a max and uses a channel to insure the handler is used concurrently by only max number of routines
func CreateThrottledHandler(handler http.Handler, maxConcurrentCalls int) ThrottledHandler {
	return ThrottledHandler{handler, make(chan struct{}, maxConcurrentCalls)}
}

// ThrottledHandler ...
type ThrottledHandler struct {
	Handler http.Handler
	ReqChan chan struct{}
}

func (m ThrottledHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) >= 3 {
		version, err := strconv.ParseFloat(pathParts[2], 64)
		if err == nil && version >= 2 { // do not default to Perl for versions over 2.x
			api.WriteRespAlertNotFound(w, r)
			return
		}
	}

	m.ReqChan <- struct{}{}
	defer func() { <-m.ReqChan }()
	m.Handler.ServeHTTP(w, r)
}
