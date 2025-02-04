// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package validation

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	envoytypev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/hashicorp/go-multierror"

	meshconfig "istio.io/api/mesh/v1alpha1"
)

func validateExtensionProviderService(service string) error {
	if service == "" {
		return fmt.Errorf("service must not be empty")
	}
	parts := strings.Split(service, "/")
	if len(parts) == 1 {
		if err := ValidateFQDN(service); err != nil {
			if err2 := ValidateIPAddress(service); err2 != nil {
				return fmt.Errorf("invalid service fmt %s: %s", service, err2)
			}
		}
	} else {
		if err := validateNamespaceSlashWildcardHostname(service, false); err != nil {
			return err
		}
	}
	return nil
}

func validateExtensionProviderEnvoyExtAuthzStatusOnError(status string) error {
	if status == "" {
		return nil
	}
	code, err := strconv.ParseInt(status, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid statusOnError value %s: %v", status, err)
	}
	if _, found := envoytypev3.StatusCode_name[int32(code)]; !found {
		return fmt.Errorf("unsupported statusOnError value %s, supported values: %v", status, envoytypev3.StatusCode_name)
	}
	return nil
}

func ValidateExtensionProviderEnvoyExtAuthzHTTP(config *meshconfig.MeshConfig_ExtensionProvider_EnvoyExternalAuthorizationHttpProvider) (errs error) {
	if config == nil {
		return fmt.Errorf("nil EnvoyExternalAuthorizationHttpProvider")
	}
	if err := ValidatePort(int(config.Port)); err != nil {
		errs = appendErrors(errs, err)
	}
	if err := validateExtensionProviderService(config.Service); err != nil {
		errs = appendErrors(errs, err)
	}
	if err := validateExtensionProviderEnvoyExtAuthzStatusOnError(config.StatusOnError); err != nil {
		errs = appendErrors(errs, err)
	}
	if config.PathPrefix != "" {
		if _, err := url.Parse(config.PathPrefix); err != nil {
			errs = appendErrors(errs, fmt.Errorf("invalid pathPrefix %s: %v", config.PathPrefix, err))
		}
		if !strings.HasPrefix(config.PathPrefix, "/") {
			errs = appendErrors(errs, fmt.Errorf("pathPrefix should begin with `/` but found %q", config.PathPrefix))
		}
	}
	return
}

func ValidateExtensionProviderEnvoyExtAuthzGRPC(config *meshconfig.MeshConfig_ExtensionProvider_EnvoyExternalAuthorizationGrpcProvider) (errs error) {
	if config == nil {
		return fmt.Errorf("nil EnvoyExternalAuthorizationGrpcProvider")
	}
	if err := ValidatePort(int(config.Port)); err != nil {
		errs = appendErrors(errs, fmt.Errorf("invalid service port: %v", err))
	}
	if err := validateExtensionProviderService(config.Service); err != nil {
		errs = appendErrors(errs, err)
	}
	if err := validateExtensionProviderEnvoyExtAuthzStatusOnError(config.StatusOnError); err != nil {
		errs = appendErrors(errs, err)
	}
	return
}

func validateExtensionProviderTracingZipkin(config *meshconfig.MeshConfig_ExtensionProvider_ZipkinTracingProvider) (errs error) {
	if config == nil {
		return fmt.Errorf("nil TracingZipkinProvider")
	}
	if err := validateExtensionProviderService(config.Service); err != nil {
		errs = appendErrors(errs, err)
	}
	if err := ValidatePort(int(config.Port)); err != nil {
		errs = appendErrors(errs, fmt.Errorf("invalid service port: %v", err))
	}
	return
}

func validateExtensionProviderTracingLightStep(config *meshconfig.MeshConfig_ExtensionProvider_LightstepTracingProvider) (errs error) {
	if config == nil {
		return fmt.Errorf("nil TracingLightStepProvider")
	}
	if err := validateExtensionProviderService(config.Service); err != nil {
		errs = appendErrors(errs, err)
	}
	if err := ValidatePort(int(config.Port)); err != nil {
		errs = appendErrors(errs, fmt.Errorf("invalid service port: %v", err))
	}
	if config.AccessToken == "" {
		errs = appendErrors(errs, fmt.Errorf("access token is required"))
	}
	return
}

func validateExtensionProviderTracingDatadog(config *meshconfig.MeshConfig_ExtensionProvider_DatadogTracingProvider) (errs error) {
	if config == nil {
		return fmt.Errorf("nil TracingDatadogProvider")
	}
	if err := validateExtensionProviderService(config.Service); err != nil {
		errs = appendErrors(errs, err)
	}
	if err := ValidatePort(int(config.Port)); err != nil {
		errs = appendErrors(errs, fmt.Errorf("invalid service port: %v", err))
	}
	return
}

func validateExtensionProviderTracingOpenCensusAgent(config *meshconfig.MeshConfig_ExtensionProvider_OpenCensusAgentTracingProvider) (errs error) {
	if config == nil {
		return fmt.Errorf("nil OpenCensusAgent")
	}
	if err := validateExtensionProviderService(config.Service); err != nil {
		errs = appendErrors(errs, err)
	}
	if err := ValidatePort(int(config.Port)); err != nil {
		errs = appendErrors(errs, fmt.Errorf("invalid service port: %v", err))
	}
	return
}

func validateExtensionProviderTracingSkyWalking(config *meshconfig.MeshConfig_ExtensionProvider_SkyWalkingTracingProvider) (errs error) {
	if config == nil {
		return fmt.Errorf("nil TracingSkyWalkingProvider")
	}
	if err := validateExtensionProviderService(config.Service); err != nil {
		errs = appendErrors(errs, err)
	}
	if err := ValidatePort(int(config.Port)); err != nil {
		errs = appendErrors(errs, fmt.Errorf("invalid service port: %v", err))
	}
	return
}

func validateExtensionProviderMetricsPrometheus(_ *meshconfig.MeshConfig_ExtensionProvider_PrometheusMetricsProvider) error {
	return nil
}

func validateExtensionProviderStackdriver(_ *meshconfig.MeshConfig_ExtensionProvider_StackdriverProvider) error {
	return nil
}

func validateExtensionProviderEnvoyFileAccessLog(_ *meshconfig.MeshConfig_ExtensionProvider_EnvoyFileAccessLogProvider) error {
	return nil
}

func ValidateExtensionProviderEnvoyOtelAls(provider *meshconfig.MeshConfig_ExtensionProvider_EnvoyOpenTelemetryLogProvider) (errs error) {
	if provider == nil {
		return fmt.Errorf("nil EnvoyOpenTelemetryLogProvider")
	}
	if err := ValidatePort(int(provider.Port)); err != nil {
		errs = appendErrors(errs, err)
	}
	if err := validateExtensionProviderService(provider.Service); err != nil {
		errs = appendErrors(errs, err)
	}
	return
}

func ValidateExtensionProviderTracingOpentelemetry(provider *meshconfig.MeshConfig_ExtensionProvider_OpenTelemetryTracingProvider) (errs error) {
	if provider == nil {
		return fmt.Errorf("nil OpenTelemetryTracingProvider")
	}
	if err := ValidatePort(int(provider.Port)); err != nil {
		errs = appendErrors(errs, err)
	}
	if err := validateExtensionProviderService(provider.Service); err != nil {
		errs = appendErrors(errs, err)
	}
	return
}

func ValidateExtensionProviderEnvoyHTTPAls(provider *meshconfig.MeshConfig_ExtensionProvider_EnvoyHttpGrpcV3LogProvider) (errs error) {
	if provider == nil {
		return fmt.Errorf("nil EnvoyHttpGrpcV3LogProvider")
	}
	if err := ValidatePort(int(provider.Port)); err != nil {
		errs = appendErrors(errs, err)
	}
	if err := validateExtensionProviderService(provider.Service); err != nil {
		errs = appendErrors(errs, err)
	}
	return
}

func ValidateExtensionProviderEnvoyTCPAls(provider *meshconfig.MeshConfig_ExtensionProvider_EnvoyTcpGrpcV3LogProvider) (errs error) {
	if provider == nil {
		return fmt.Errorf("nil EnvoyTcpGrpcV3LogProvider")
	}
	if err := ValidatePort(int(provider.Port)); err != nil {
		errs = appendErrors(errs, err)
	}
	if err := validateExtensionProviderService(provider.Service); err != nil {
		errs = appendErrors(errs, err)
	}
	return
}

func validateExtensionProvider(config *meshconfig.MeshConfig) (errs error) {
	definedProviders := map[string]struct{}{}
	for _, c := range config.ExtensionProviders {
		var currentErrs error
		// Provider name must be unique and not empty.
		if c.Name == "" {
			currentErrs = appendErrors(currentErrs, fmt.Errorf("empty extension provider name"))
		} else {
			if _, found := definedProviders[c.Name]; found {
				currentErrs = appendErrors(currentErrs, fmt.Errorf("duplicate extension provider name %s", c.Name))
			}
			definedProviders[c.Name] = struct{}{}
		}

		switch provider := c.Provider.(type) {
		case *meshconfig.MeshConfig_ExtensionProvider_EnvoyExtAuthzHttp:
			currentErrs = appendErrors(currentErrs, ValidateExtensionProviderEnvoyExtAuthzHTTP(provider.EnvoyExtAuthzHttp))
		case *meshconfig.MeshConfig_ExtensionProvider_EnvoyExtAuthzGrpc:
			currentErrs = appendErrors(currentErrs, ValidateExtensionProviderEnvoyExtAuthzGRPC(provider.EnvoyExtAuthzGrpc))
		case *meshconfig.MeshConfig_ExtensionProvider_Zipkin:
			currentErrs = appendErrors(currentErrs, validateExtensionProviderTracingZipkin(provider.Zipkin))
		case *meshconfig.MeshConfig_ExtensionProvider_Lightstep:
			currentErrs = appendErrors(currentErrs, validateExtensionProviderTracingLightStep(provider.Lightstep))
		case *meshconfig.MeshConfig_ExtensionProvider_Datadog:
			currentErrs = appendErrors(currentErrs, validateExtensionProviderTracingDatadog(provider.Datadog))
		case *meshconfig.MeshConfig_ExtensionProvider_Opencensus:
			currentErrs = appendErrors(currentErrs, validateExtensionProviderTracingOpenCensusAgent(provider.Opencensus))
		case *meshconfig.MeshConfig_ExtensionProvider_Skywalking:
			currentErrs = appendErrors(currentErrs, validateExtensionProviderTracingSkyWalking(provider.Skywalking))
		case *meshconfig.MeshConfig_ExtensionProvider_Prometheus:
			currentErrs = appendErrors(currentErrs, validateExtensionProviderMetricsPrometheus(provider.Prometheus))
		case *meshconfig.MeshConfig_ExtensionProvider_Stackdriver:
			currentErrs = appendErrors(currentErrs, validateExtensionProviderStackdriver(provider.Stackdriver))
		case *meshconfig.MeshConfig_ExtensionProvider_EnvoyFileAccessLog:
			currentErrs = appendErrors(currentErrs, validateExtensionProviderEnvoyFileAccessLog(provider.EnvoyFileAccessLog))
		case *meshconfig.MeshConfig_ExtensionProvider_EnvoyOtelAls:
			currentErrs = appendErrors(currentErrs, ValidateExtensionProviderEnvoyOtelAls(provider.EnvoyOtelAls))
		case *meshconfig.MeshConfig_ExtensionProvider_Opentelemetry:
			currentErrs = appendErrors(currentErrs, ValidateExtensionProviderTracingOpentelemetry(provider.Opentelemetry))
		case *meshconfig.MeshConfig_ExtensionProvider_EnvoyHttpAls:
			currentErrs = appendErrors(currentErrs, ValidateExtensionProviderEnvoyHTTPAls(provider.EnvoyHttpAls))
		case *meshconfig.MeshConfig_ExtensionProvider_EnvoyTcpAls:
			currentErrs = appendErrors(currentErrs, ValidateExtensionProviderEnvoyTCPAls(provider.EnvoyTcpAls))
			// TODO: add exhaustiveness test
		default:
			currentErrs = appendErrors(currentErrs, fmt.Errorf("unsupported provider: %v of type %T", provider, provider))
		}
		currentErrs = multierror.Prefix(currentErrs, fmt.Sprintf("invalid extension provider %s:", c.Name))
		errs = appendErrors(errs, currentErrs)
	}
	return
}
