package paas

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// InfraIntegration provides integration with the existing infra codebase
type InfraIntegration struct {
	InfraPath string // Path to the infra directory
}

// NewInfraIntegration creates a new infra integration instance
func NewInfraIntegration(infraPath string) *InfraIntegration {
	return &InfraIntegration{
		InfraPath: infraPath,
	}
}

// DeployToInfra converts an Application to infra Go code and writes service files
func (ii *InfraIntegration) DeployToInfra(app *Application, servicesDir string) error {
	servicesPath := filepath.Join(ii.InfraPath, servicesDir)
	if err := os.MkdirAll(servicesPath, 0755); err != nil {
		return fmt.Errorf("failed to create services directory: %w", err)
	}

	// Generate service files
	for name, service := range app.Services {
		filename := fmt.Sprintf("services_%s.go", strings.ToLower(name))
		filepath := filepath.Join(servicesPath, filename)

		code, err := ii.generateInfraService(name, service)
		if err != nil {
			return fmt.Errorf("failed to generate service code for %s: %w", name, err)
		}

		if err := os.WriteFile(filepath, []byte(code), 0644); err != nil {
			return fmt.Errorf("failed to write service file %s: %w", filepath, err)
		}

		fmt.Printf("Generated %s\n", filepath)
	}

	// Update main services.go to include new services
	if err := ii.updateMainServicesFile(app); err != nil {
		return fmt.Errorf("failed to update main services file: %w", err)
	}

	return nil
}

// generateInfraService generates Go code for a service in infra format
func (ii *InfraIntegration) generateInfraService(name string, service *Service) (string, error) {
	var sb strings.Builder

	sb.WriteString("package main\n\n")
	sb.WriteString("import (\n")
	sb.WriteString("\t\"fmt\"\n")
	sb.WriteString("\t\"github.com/docker/docker/api/types/mount\"\n")
	sb.WriteString(")\n\n")

	sb.WriteString(fmt.Sprintf("// %sService returns the %s service configuration\n", strings.Title(name), name))
	sb.WriteString(fmt.Sprintf("func %sService(configPath, secretsPath, domain, tsHostname string) Service {\n", strings.Title(name)))

	// Basic service fields
	sb.WriteString(fmt.Sprintf("\treturn Service{\n"))
	sb.WriteString(fmt.Sprintf("\t\tName:          \"%s\",\n", name))
	sb.WriteString(fmt.Sprintf("\t\tImage:         \"%s\",\n", service.Image))

	if service.ContainerName != "" {
		sb.WriteString(fmt.Sprintf("\t\tContainerName: \"%s\",\n", service.ContainerName))
	}

	if service.Hostname != "" {
		sb.WriteString(fmt.Sprintf("\t\tHostname:      \"%s\",\n", service.Hostname))
	}

	// Ports
	if len(service.Ports) > 0 {
		sb.WriteString("\t\tPorts: []PortMapping{\n")
		for _, port := range service.Ports {
			sb.WriteString("\t\t\t{")
			if port.Name != "" {
				sb.WriteString(fmt.Sprintf("Name: \"%s\", ", port.Name))
			}
			if port.TargetName != "" {
				sb.WriteString(fmt.Sprintf("TargetName: \"%s\", ", port.TargetName))
			}
			if port.HostIP != "" {
				sb.WriteString(fmt.Sprintf("HostIP: \"%s\", ", port.HostIP))
			}
			sb.WriteString(fmt.Sprintf("HostPort: \"%s\", ContainerPort: \"%s\"", port.HostPort, port.ContainerPort))
			if port.Protocol != "" {
				sb.WriteString(fmt.Sprintf(", Protocol: \"%s\"", port.Protocol))
			}
			if len(port.Extensions) > 0 {
				sb.WriteString(fmt.Sprintf(", Extensions: %#v", port.Extensions))
			}
			sb.WriteString("},\n")
		}
		sb.WriteString("\t\t},\n")
	}

	// Environment variables
	if len(service.Environment) > 0 {
		sb.WriteString("\t\tEnvironment: map[string]string{\n")
		for key, value := range service.Environment {
			sb.WriteString(fmt.Sprintf("\t\t\t\"%s\": \"%s\",\n", key, value))
		}
		sb.WriteString("\t\t},\n")
	}

	// Volumes
	if len(service.Volumes) > 0 {
		sb.WriteString("\t\tVolumes: []VolumeMount{\n")
		for _, vol := range service.Volumes {
			sb.WriteString("\t\t\t{")
			if vol.Source != "" {
				sb.WriteString(fmt.Sprintf("Source: \"%s\", ", vol.Source))
			}
			if vol.Target != "" {
				sb.WriteString(fmt.Sprintf("Target: \"%s\", ", vol.Target))
			}
			if vol.Type != "" {
				sb.WriteString(fmt.Sprintf("Type: \"%s\", ", vol.Type))
			}
			if vol.ReadOnly {
				sb.WriteString("ReadOnly: true, ")
			}
			if vol.Mode != "" {
				sb.WriteString(fmt.Sprintf("Mode: \"%s\", ", vol.Mode))
			}
			if vol.Propagation != "" {
				sb.WriteString(fmt.Sprintf("Propagation: \"%s\", ", vol.Propagation))
			}
			if vol.NoCopy {
				sb.WriteString("NoCopy: true, ")
			}
			if len(vol.Options) > 0 {
				sb.WriteString("Options: map[string]string{")
				for _, key := range sortedMapKeys(vol.Options) {
					sb.WriteString(fmt.Sprintf("\"%s\": \"%s\", ", key, vol.Options[key]))
				}
				sb.WriteString("}, ")
			}
			sb.WriteString("},\n")
		}
		sb.WriteString("\t\t},\n")
	}

	if len(service.Configs) > 0 {
		sb.WriteString("\t\tConfigs: []ConfigMount{\n")
		for _, cfg := range service.Configs {
			sb.WriteString("\t\t\t{")
			if cfg.Source != "" {
				sb.WriteString(fmt.Sprintf("Source: \"%s\", ", cfg.Source))
			}
			if cfg.Target != "" {
				sb.WriteString(fmt.Sprintf("Target: \"%s\", ", cfg.Target))
			}
			if cfg.Mode != "" {
				sb.WriteString(fmt.Sprintf("Mode: \"%s\", ", cfg.Mode))
			}
			sb.WriteString("},\n")
		}
		sb.WriteString("\t\t},\n")
	}

	if len(service.Secrets) > 0 {
		sb.WriteString("\t\tSecrets: []SecretMount{\n")
		for _, secret := range service.Secrets {
			sb.WriteString("\t\t\t{")
			if secret.Source != "" {
				sb.WriteString(fmt.Sprintf("Source: \"%s\", ", secret.Source))
			}
			if secret.Target != "" {
				sb.WriteString(fmt.Sprintf("Target: \"%s\", ", secret.Target))
			}
			if secret.Mode != "" {
				sb.WriteString(fmt.Sprintf("Mode: \"%s\", ", secret.Mode))
			}
			sb.WriteString("},\n")
		}
		sb.WriteString("\t\t},\n")
	}

	if len(service.DNS) > 0 {
		sb.WriteString(fmt.Sprintf("\t\tDNS: []string{%s},\n", ii.formatStringSlice(service.DNS)))
	}
	if len(service.DNSSearch) > 0 {
		sb.WriteString(fmt.Sprintf("\t\tDNSSearch: []string{%s},\n", ii.formatStringSlice(service.DNSSearch)))
	}
	if len(service.DNSOptions) > 0 {
		sb.WriteString(fmt.Sprintf("\t\tDNSOptions: []string{%s},\n", ii.formatStringSlice(service.DNSOptions)))
	}
	if len(service.ExtraHosts) > 0 {
		sb.WriteString(fmt.Sprintf("\t\tExtraHosts: []string{%s},\n", ii.formatStringSlice(service.ExtraHosts)))
	}

	// Networks
	if len(service.Networks) > 0 {
		sb.WriteString(fmt.Sprintf("\t\tNetworks:      []string{%s},\n",
			ii.formatStringSlice(service.Networks)))
	}

	if len(service.Command) > 0 {
		sb.WriteString(fmt.Sprintf("\t\tCommand: []string{%s},\n", ii.formatStringSlice(service.Command)))
	}
	if len(service.Entrypoint) > 0 {
		sb.WriteString(fmt.Sprintf("\t\tEntrypoint: []string{%s},\n", ii.formatStringSlice(service.Entrypoint)))
	}
	if service.WorkingDir != "" {
		sb.WriteString(fmt.Sprintf("\t\tWorkingDir: \"%s\",\n", service.WorkingDir))
	}
	if service.Build != nil {
		sb.WriteString("\t\tBuild: &BuildConfig{\n")
		if service.Build.Context != "" {
			sb.WriteString(fmt.Sprintf("\t\t\tContext: \"%s\",\n", service.Build.Context))
		}
		if service.Build.Dockerfile != "" {
			sb.WriteString(fmt.Sprintf("\t\t\tDockerfile: \"%s\",\n", service.Build.Dockerfile))
		}
		if len(service.Build.Extensions) > 0 {
			sb.WriteString(fmt.Sprintf("\t\t\tExtensions: %#v,\n", service.Build.Extensions))
		}
		sb.WriteString("\t\t},\n")
	}
	if len(service.Devices) > 0 {
		sb.WriteString(fmt.Sprintf("\t\tDevices: []string{%s},\n", ii.formatStringSlice(service.Devices)))
	}
	if len(service.Expose) > 0 {
		sb.WriteString(fmt.Sprintf("\t\tExpose: []string{%s},\n", ii.formatStringSlice(service.Expose)))
	}

	// Restart policy
	if service.Restart != "" {
		sb.WriteString(fmt.Sprintf("\t\tRestart:       \"%s\",\n", service.Restart))
	}

	if service.Privileged {
		sb.WriteString("\t\tPrivileged: true,\n")
	}
	if service.User != "" {
		sb.WriteString(fmt.Sprintf("\t\tUser: \"%s\",\n", service.User))
	}
	if service.Group != "" {
		sb.WriteString(fmt.Sprintf("\t\tGroup: \"%s\",\n", service.Group))
	}
	if len(service.CapAdd) > 0 {
		sb.WriteString(fmt.Sprintf("\t\tCapAdd: []string{%s},\n", ii.formatStringSlice(service.CapAdd)))
	}
	if len(service.CapDrop) > 0 {
		sb.WriteString(fmt.Sprintf("\t\tCapDrop: []string{%s},\n", ii.formatStringSlice(service.CapDrop)))
	}
	if len(service.SecurityOpt) > 0 {
		sb.WriteString(fmt.Sprintf("\t\tSecurityOpt: []string{%s},\n", ii.formatStringSlice(service.SecurityOpt)))
	}
	if service.ReadOnlyRootFS {
		sb.WriteString("\t\tReadOnlyRootFS: true,\n")
	}
	if service.Init != nil {
		if *service.Init {
			sb.WriteString("\t\tInit: func() *bool { v := true; return &v }(),\n")
		} else {
			sb.WriteString("\t\tInit: func() *bool { v := false; return &v }(),\n")
		}
	}
	if service.Tty {
		sb.WriteString("\t\tTty: true,\n")
	}
	if service.StdinOpen {
		sb.WriteString("\t\tStdinOpen: true,\n")
	}
	if service.StopSignal != "" {
		sb.WriteString(fmt.Sprintf("\t\tStopSignal: \"%s\",\n", service.StopSignal))
	}
	if service.StopGracePeriod != "" {
		sb.WriteString(fmt.Sprintf("\t\tStopGracePeriod: \"%s\",\n", service.StopGracePeriod))
	}
	if service.Replicas > 0 {
		sb.WriteString(fmt.Sprintf("\t\tReplicas: %d,\n", service.Replicas))
	}
	if service.MemLimit != "" {
		sb.WriteString(fmt.Sprintf("\t\tMemLimit: \"%s\",\n", service.MemLimit))
	}
	if service.MemReservation != "" {
		sb.WriteString(fmt.Sprintf("\t\tMemReservation: \"%s\",\n", service.MemReservation))
	}
	if service.CPUs != "" {
		sb.WriteString(fmt.Sprintf("\t\tCPUs: \"%s\",\n", service.CPUs))
	}
	if service.UserNSMode != "" {
		sb.WriteString(fmt.Sprintf("\t\tUserNSMode: \"%s\",\n", service.UserNSMode))
	}
	if service.PullPolicy != "" {
		sb.WriteString(fmt.Sprintf("\t\tPullPolicy: \"%s\",\n", service.PullPolicy))
	}
	if service.Ulimits != nil && service.Ulimits.Nofile != nil {
		sb.WriteString("\t\tUlimits: &Ulimits{\n")
		sb.WriteString("\t\t\tNofile: &NofileLimit{\n")
		sb.WriteString(fmt.Sprintf("\t\t\t\tSoft: %d,\n", service.Ulimits.Nofile.Soft))
		sb.WriteString(fmt.Sprintf("\t\t\t\tHard: %d,\n", service.Ulimits.Nofile.Hard))
		sb.WriteString("\t\t\t},\n")
		sb.WriteString("\t\t},\n")
	}
	if service.HealthCheck != nil {
		sb.WriteString("\t\tHealthcheck: &Healthcheck{\n")
		if len(service.HealthCheck.Test) > 0 {
			sb.WriteString(fmt.Sprintf("\t\t\tTest: []string{%s},\n", ii.formatStringSlice(service.HealthCheck.Test)))
		}
		if service.HealthCheck.Interval != "" {
			sb.WriteString(fmt.Sprintf("\t\t\tInterval: \"%s\",\n", service.HealthCheck.Interval))
		}
		if service.HealthCheck.Timeout != "" {
			sb.WriteString(fmt.Sprintf("\t\t\tTimeout: \"%s\",\n", service.HealthCheck.Timeout))
		}
		if service.HealthCheck.Retries > 0 {
			sb.WriteString(fmt.Sprintf("\t\t\tRetries: %d,\n", service.HealthCheck.Retries))
		}
		if service.HealthCheck.StartPeriod != "" {
			sb.WriteString(fmt.Sprintf("\t\t\tStartPeriod: \"%s\",\n", service.HealthCheck.StartPeriod))
		}
		sb.WriteString("\t\t},\n")
	}

	// Labels
	if len(service.Labels) > 0 {
		sb.WriteString("\t\tLabels: map[string]string{\n")
		for key, value := range service.Labels {
			sb.WriteString(fmt.Sprintf("\t\t\t\"%s\": \"%s\",\n", key, value))
		}
		sb.WriteString("\t\t},\n")
	}

	sb.WriteString("\t}\n")
	sb.WriteString("}\n")

	return sb.String(), nil
}

// updateMainServicesFile updates the main services.go file to include new services
func (ii *InfraIntegration) updateMainServicesFile(app *Application) error {
	mainServicesFile := filepath.Join(ii.InfraPath, "services.go")

	// Read existing file
	content, err := os.ReadFile(mainServicesFile)
	if err != nil {
		return fmt.Errorf("failed to read main services file: %w", err)
	}

	// Parse the Go file
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse Go file: %w", err)
	}

	// Find the defineServicesFromConfig function
	var defineServicesFunc *ast.FuncDecl
	ast.Inspect(file, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok && fn.Name.Name == "defineServicesFromConfig" {
			defineServicesFunc = fn
			return false
		}
		return true
	})

	if defineServicesFunc == nil {
		return fmt.Errorf("could not find defineServicesFromConfig function")
	}

	if strings.Contains(string(content), "// PaaS-generated services") {
		return nil
	}

	insertAt := -1
	for _, stmt := range defineServicesFunc.Body.List {
		ret, ok := stmt.(*ast.ReturnStmt)
		if !ok || len(ret.Results) != 1 {
			continue
		}
		if exprString(ret.Results[0]) == "services" {
			insertAt = fset.Position(ret.Pos()).Offset
			break
		}
	}
	if insertAt < 0 {
		insertAt = fset.Position(defineServicesFunc.Body.Rbrace).Offset
	}

	var snippet strings.Builder
	snippet.WriteString("\n\t// PaaS-generated services\n")
	for name := range app.Services {
		titleName := strings.Title(name)
		snippet.WriteString(fmt.Sprintf("\tservices = append(services, %sService(configPath, secretsPath, domain, getEnv(\"TS_HOSTNAME\", \"localhost\")))\n", titleName))
	}

	newContent := string(content[:insertAt]) + snippet.String() + string(content[insertAt:])

	// Format the code
	formatted, err := format.Source([]byte(newContent))
	if err != nil {
		return fmt.Errorf("failed to format Go code: %w", err)
	}

	// Write back
	return os.WriteFile(mainServicesFile, formatted, 0644)
}

// LoadFromInfra loads an Application from existing infra Go files
func (ii *InfraIntegration) LoadFromInfra() (*Application, error) {
	servicesPath := filepath.Join(ii.InfraPath, "services_generated")
	if _, err := os.Stat(servicesPath); err != nil {
		return nil, fmt.Errorf("infra services directory not available: %w", err)
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, servicesPath, func(info os.FileInfo) bool {
		return !info.IsDir() && strings.HasSuffix(info.Name(), ".go")
	}, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse infra services: %w", err)
	}

	app := &Application{
		Platform: PlatformDockerCompose,
		Services: make(map[string]*Service),
	}

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok || fn.Body == nil || fn.Type.Results == nil || len(fn.Type.Results.List) != 1 {
					continue
				}
				if fn.Type.Results.List[0].Type == nil || exprString(fn.Type.Results.List[0].Type) != "Service" {
					continue
				}
				service := parseInfraServiceFunc(fn)
				if service == nil {
					continue
				}
				if service.Name == "" {
					service.Name = strings.TrimSuffix(fn.Name.Name, "Service")
					service.Name = strings.ToLower(service.Name)
				}
				if service.Platform == "" {
					service.Platform = PlatformDockerCompose
				}
				app.Services[service.Name] = service
			}
		}
	}

	if len(app.Services) == 0 {
		return app, fmt.Errorf("no infra services discovered in %s", servicesPath)
	}

	app.AttachCanonical()
	return app, nil
}

func parseInfraServiceFunc(fn *ast.FuncDecl) *Service {
	for _, stmt := range fn.Body.List {
		ret, ok := stmt.(*ast.ReturnStmt)
		if !ok || len(ret.Results) == 0 {
			continue
		}
		lit, ok := ret.Results[0].(*ast.CompositeLit)
		if !ok || exprString(lit.Type) != "Service" {
			continue
		}
		return parseInfraServiceLiteral(lit)
	}
	return nil
}

func parseInfraServiceLiteral(lit *ast.CompositeLit) *Service {
	service := &Service{
		Environment: map[string]string{},
		Labels:      map[string]string{},
		Extensions:  map[string]interface{}{},
	}
	for _, element := range lit.Elts {
		field, ok := element.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key := exprString(field.Key)
		switch key {
		case "Name":
			service.Name = stringLiteral(field.Value)
		case "Image":
			service.Image = stringLiteral(field.Value)
		case "ContainerName":
			service.ContainerName = stringLiteral(field.Value)
		case "Hostname":
			service.Hostname = stringLiteral(field.Value)
		case "Ports":
			service.Ports = parseInfraPortMappings(field.Value)
		case "Environment":
			service.Environment = parseInfraStringMap(field.Value)
		case "DNS":
			service.DNS = parseInfraStringSlice(field.Value)
		case "DNSSearch":
			service.DNSSearch = parseInfraStringSlice(field.Value)
		case "DNSOptions":
			service.DNSOptions = parseInfraStringSlice(field.Value)
		case "ExtraHosts":
			service.ExtraHosts = parseInfraStringSlice(field.Value)
		case "Volumes":
			service.Volumes = parseInfraVolumeMounts(field.Value)
		case "Configs":
			service.Configs = parseInfraConfigMounts(field.Value)
		case "Secrets":
			service.Secrets = parseInfraSecretMounts(field.Value)
		case "Networks":
			service.Networks = parseInfraStringSlice(field.Value)
		case "Command":
			service.Command = parseInfraStringSlice(field.Value)
		case "Entrypoint":
			service.Entrypoint = parseInfraStringSlice(field.Value)
		case "WorkingDir":
			service.WorkingDir = stringLiteral(field.Value)
		case "Build":
			service.Build = parseInfraBuild(field.Value)
		case "Devices":
			service.Devices = parseInfraStringSlice(field.Value)
		case "Expose":
			service.Expose = parseInfraStringSlice(field.Value)
		case "Restart":
			service.Restart = stringLiteral(field.Value)
		case "Labels":
			service.Labels = parseInfraStringMap(field.Value)
		case "Privileged":
			service.Privileged = boolLiteral(field.Value)
		case "User":
			service.User = stringLiteral(field.Value)
		case "Group":
			service.Group = stringLiteral(field.Value)
		case "CapAdd":
			service.CapAdd = parseInfraStringSlice(field.Value)
		case "CapDrop":
			service.CapDrop = parseInfraStringSlice(field.Value)
		case "SecurityOpt":
			service.SecurityOpt = parseInfraStringSlice(field.Value)
		case "ReadOnlyRootFS":
			service.ReadOnlyRootFS = boolLiteral(field.Value)
		case "Init":
			service.Init = parseInfraBoolPtr(field.Value)
		case "Tty":
			service.Tty = boolLiteral(field.Value)
		case "StdinOpen":
			service.StdinOpen = boolLiteral(field.Value)
		case "StopSignal":
			service.StopSignal = stringLiteral(field.Value)
		case "StopGracePeriod":
			service.StopGracePeriod = stringLiteral(field.Value)
		case "Replicas":
			service.Replicas = intLiteral(field.Value)
		case "MemLimit":
			service.MemLimit = stringLiteral(field.Value)
		case "MemReservation":
			service.MemReservation = stringLiteral(field.Value)
		case "CPUs":
			service.CPUs = stringLiteral(field.Value)
		case "UserNSMode":
			service.UserNSMode = stringLiteral(field.Value)
		case "PullPolicy":
			service.PullPolicy = stringLiteral(field.Value)
		case "Ulimits":
			service.Ulimits = parseInfraUlimits(field.Value)
		case "Healthcheck":
			service.HealthCheck = parseInfraHealthcheck(field.Value)
		}
	}
	return service
}

func parseInfraBoolPtr(expr ast.Expr) *bool {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return nil
	}
	lit, ok := call.Fun.(*ast.FuncLit)
	if !ok || lit.Body == nil {
		return nil
	}
	for _, stmt := range lit.Body.List {
		assign, ok := stmt.(*ast.AssignStmt)
		if !ok || len(assign.Rhs) != 1 {
			continue
		}
		if ident, ok := assign.Rhs[0].(*ast.Ident); ok && (ident.Name == "true" || ident.Name == "false") {
			value := ident.Name == "true"
			return &value
		}
	}
	return nil
}

func parseInfraPortMappings(expr ast.Expr) []PortMapping {
	lit, ok := expr.(*ast.CompositeLit)
	if !ok {
		return nil
	}
	var ports []PortMapping
	for _, element := range lit.Elts {
		portLit, ok := element.(*ast.CompositeLit)
		if !ok {
			continue
		}
		port := PortMapping{}
		for _, fieldExpr := range portLit.Elts {
			field, ok := fieldExpr.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			switch exprString(field.Key) {
			case "Name":
				port.Name = stringLiteral(field.Value)
			case "TargetName":
				port.TargetName = stringLiteral(field.Value)
			case "HostIP":
				port.HostIP = stringLiteral(field.Value)
			case "HostPort":
				port.HostPort = stringLiteral(field.Value)
			case "ContainerPort":
				port.ContainerPort = stringLiteral(field.Value)
			case "Protocol":
				port.Protocol = stringLiteral(field.Value)
			case "Extensions":
				port.Extensions = parseInfraInterfaceMap(field.Value)
			}
		}
		if len(port.Extensions) == 0 {
			port.Extensions = nil
		}
		ports = append(ports, port)
	}
	return ports
}

func parseInfraVolumeMounts(expr ast.Expr) []VolumeMount {
	lit, ok := expr.(*ast.CompositeLit)
	if !ok {
		return nil
	}
	var mounts []VolumeMount
	for _, element := range lit.Elts {
		mountLit, ok := element.(*ast.CompositeLit)
		if !ok {
			continue
		}
		mount := VolumeMount{Options: map[string]string{}}
		for _, fieldExpr := range mountLit.Elts {
			field, ok := fieldExpr.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			switch exprString(field.Key) {
			case "Source":
				mount.Source = stringLiteral(field.Value)
			case "Target":
				mount.Target = stringLiteral(field.Value)
			case "Type":
				mount.Type = stringLiteral(field.Value)
			case "ReadOnly":
				mount.ReadOnly = boolLiteral(field.Value)
			case "Mode":
				mount.Mode = stringLiteral(field.Value)
			case "Propagation":
				mount.Propagation = stringLiteral(field.Value)
			case "NoCopy":
				mount.NoCopy = boolLiteral(field.Value)
			case "Options":
				mount.Options = parseInfraStringMap(field.Value)
			}
		}
		mounts = append(mounts, mount)
	}
	return mounts
}

func parseInfraConfigMounts(expr ast.Expr) []FileRef {
	lit, ok := expr.(*ast.CompositeLit)
	if !ok {
		return nil
	}
	var mounts []FileRef
	for _, element := range lit.Elts {
		mountLit, ok := element.(*ast.CompositeLit)
		if !ok {
			continue
		}
		mount := FileRef{}
		for _, fieldExpr := range mountLit.Elts {
			field, ok := fieldExpr.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			switch exprString(field.Key) {
			case "Source":
				mount.Source = stringLiteral(field.Value)
			case "Target":
				mount.Target = stringLiteral(field.Value)
			case "Mode":
				mount.Mode = stringLiteral(field.Value)
			}
		}
		mounts = append(mounts, mount)
	}
	return mounts
}

func parseInfraSecretMounts(expr ast.Expr) []FileRef {
	lit, ok := expr.(*ast.CompositeLit)
	if !ok {
		return nil
	}
	var mounts []FileRef
	for _, element := range lit.Elts {
		mountLit, ok := element.(*ast.CompositeLit)
		if !ok {
			continue
		}
		mount := FileRef{}
		for _, fieldExpr := range mountLit.Elts {
			field, ok := fieldExpr.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			switch exprString(field.Key) {
			case "Source":
				mount.Source = stringLiteral(field.Value)
			case "Target":
				mount.Target = stringLiteral(field.Value)
			case "Mode":
				mount.Mode = stringLiteral(field.Value)
			}
		}
		mounts = append(mounts, mount)
	}
	return mounts
}

func parseInfraHealthcheck(expr ast.Expr) *HealthCheck {
	switch value := expr.(type) {
	case *ast.UnaryExpr:
		if value.Op == token.AND {
			if lit, ok := value.X.(*ast.CompositeLit); ok {
				return parseInfraHealthcheckLiteral(lit)
			}
		}
	case *ast.CompositeLit:
		return parseInfraHealthcheckLiteral(value)
	}
	return nil
}

func parseInfraHealthcheckLiteral(lit *ast.CompositeLit) *HealthCheck {
	health := &HealthCheck{}
	for _, element := range lit.Elts {
		field, ok := element.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		switch exprString(field.Key) {
		case "Test":
			health.Test = parseInfraStringSlice(field.Value)
		case "Interval":
			health.Interval = stringLiteral(field.Value)
		case "Timeout":
			health.Timeout = stringLiteral(field.Value)
		case "Retries":
			health.Retries = intLiteral(field.Value)
		case "StartPeriod":
			health.StartPeriod = stringLiteral(field.Value)
		}
	}
	return health
}

func parseInfraBuild(expr ast.Expr) *BuildConfig {
	switch value := expr.(type) {
	case *ast.UnaryExpr:
		if value.Op == token.AND {
			if lit, ok := value.X.(*ast.CompositeLit); ok {
				return parseInfraBuildLiteral(lit)
			}
		}
	case *ast.CompositeLit:
		return parseInfraBuildLiteral(value)
	}
	return nil
}

func parseInfraBuildLiteral(lit *ast.CompositeLit) *BuildConfig {
	build := &BuildConfig{Extensions: map[string]interface{}{}}
	for _, element := range lit.Elts {
		field, ok := element.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		switch exprString(field.Key) {
		case "Context":
			build.Context = stringLiteral(field.Value)
		case "Dockerfile":
			build.Dockerfile = stringLiteral(field.Value)
		case "Extensions":
			build.Extensions = parseInfraInterfaceMap(field.Value)
		}
	}
	if len(build.Extensions) == 0 {
		build.Extensions = nil
	}
	return build
}

func parseInfraUlimits(expr ast.Expr) *Ulimits {
	switch value := expr.(type) {
	case *ast.UnaryExpr:
		if value.Op == token.AND {
			if lit, ok := value.X.(*ast.CompositeLit); ok {
				return parseInfraUlimitsLiteral(lit)
			}
		}
	case *ast.CompositeLit:
		return parseInfraUlimitsLiteral(value)
	}
	return nil
}

func parseInfraUlimitsLiteral(lit *ast.CompositeLit) *Ulimits {
	limits := &Ulimits{}
	for _, element := range lit.Elts {
		field, ok := element.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		if exprString(field.Key) != "Nofile" {
			continue
		}
		switch value := field.Value.(type) {
		case *ast.UnaryExpr:
			if value.Op == token.AND {
				if composite, ok := value.X.(*ast.CompositeLit); ok {
					limits.Nofile = parseInfraNofileLimit(composite)
				}
			}
		case *ast.CompositeLit:
			limits.Nofile = parseInfraNofileLimit(value)
		}
	}
	return limits
}

func parseInfraNofileLimit(lit *ast.CompositeLit) *NofileLimit {
	limit := &NofileLimit{}
	for _, element := range lit.Elts {
		field, ok := element.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		switch exprString(field.Key) {
		case "Soft":
			limit.Soft = intLiteral(field.Value)
		case "Hard":
			limit.Hard = intLiteral(field.Value)
		}
	}
	return limit
}

func parseInfraStringSlice(expr ast.Expr) []string {
	lit, ok := expr.(*ast.CompositeLit)
	if !ok {
		return nil
	}
	values := make([]string, 0, len(lit.Elts))
	for _, element := range lit.Elts {
		if value := stringLiteral(element); value != "" {
			values = append(values, value)
		}
	}
	return values
}

func parseInfraStringMap(expr ast.Expr) map[string]string {
	lit, ok := expr.(*ast.CompositeLit)
	if !ok {
		return nil
	}
	values := map[string]string{}
	for _, element := range lit.Elts {
		field, ok := element.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key := stringLiteral(field.Key)
		if key == "" {
			continue
		}
		values[key] = stringLiteral(field.Value)
	}
	return values
}

func parseInfraInterfaceMap(expr ast.Expr) map[string]interface{} {
	lit, ok := expr.(*ast.CompositeLit)
	if !ok {
		return nil
	}
	values := map[string]interface{}{}
	for _, element := range lit.Elts {
		field, ok := element.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key := stringLiteral(field.Key)
		if key == "" {
			continue
		}
		values[key] = parseInfraInterfaceValue(field.Value)
	}
	if len(values) == 0 {
		return nil
	}
	return values
}

func parseInfraInterfaceValue(expr ast.Expr) interface{} {
	switch v := expr.(type) {
	case *ast.BasicLit:
		switch v.Kind {
		case token.STRING:
			if value, err := strconv.Unquote(v.Value); err == nil {
				return value
			}
			return strings.Trim(v.Value, "\"")
		case token.INT:
			if value, err := strconv.Atoi(v.Value); err == nil {
				return value
			}
			return v.Value
		}
	case *ast.Ident:
		switch v.Name {
		case "true":
			return true
		case "false":
			return false
		case "nil":
			return nil
		default:
			return v.Name
		}
	case *ast.CompositeLit:
		if mapped := parseInfraInterfaceMap(v); mapped != nil {
			return mapped
		}
		if values := parseInfraStringSlice(v); len(values) > 0 {
			return values
		}
	}
	return nil
}

func stringLiteral(expr ast.Expr) string {
	switch v := expr.(type) {
	case *ast.BasicLit:
		if v.Kind == token.STRING {
			value, err := strconv.Unquote(v.Value)
			if err == nil {
				return value
			}
			return strings.Trim(v.Value, "\"")
		}
	case *ast.Ident:
		if v.Name != "nil" {
			return v.Name
		}
	}
	return ""
}

func boolLiteral(expr ast.Expr) bool {
	if ident, ok := expr.(*ast.Ident); ok {
		return ident.Name == "true"
	}
	return false
}

func intLiteral(expr ast.Expr) int {
	if lit, ok := expr.(*ast.BasicLit); ok && lit.Kind == token.INT {
		if value, err := strconv.Atoi(lit.Value); err == nil {
			return value
		}
	}
	return 0
}

func exprString(expr ast.Expr) string {
	if expr == nil {
		return ""
	}
	switch v := expr.(type) {
	case *ast.Ident:
		return v.Name
	case *ast.SelectorExpr:
		if prefix := exprString(v.X); prefix != "" {
			return prefix + "." + v.Sel.Name
		}
		return v.Sel.Name
	case *ast.StarExpr:
		return exprString(v.X)
	case *ast.ArrayType:
		return "[]" + exprString(v.Elt)
	default:
		return ""
	}
}

// formatStringSlice formats a string slice for Go code
func (ii *InfraIntegration) formatStringSlice(slice []string) string {
	if len(slice) == 0 {
		return ""
	}

	quoted := make([]string, len(slice))
	for i, s := range slice {
		quoted[i] = fmt.Sprintf("\"%s\"", s)
	}

	return strings.Join(quoted, ", ")
}

// ValidateInfraIntegration validates that the infra integration will work
func (ii *InfraIntegration) ValidateInfraIntegration() error {
	// Check if infra directory exists
	if _, err := os.Stat(ii.InfraPath); os.IsNotExist(err) {
		return fmt.Errorf("infra directory does not exist: %s", ii.InfraPath)
	}

	// Check if main services.go exists
	mainServicesFile := filepath.Join(ii.InfraPath, "services.go")
	if _, err := os.Stat(mainServicesFile); os.IsNotExist(err) {
		return fmt.Errorf("main services.go file does not exist: %s", mainServicesFile)
	}

	return nil
}
