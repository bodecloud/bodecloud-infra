import { App, Chart, ChartProps } from 'cdk8s';
import { Construct } from 'constructs';

// Import infrastructure functions
import { createNamespaces } from './lib/infrastructure/namespaces';
import { createResourceQuotas } from './lib/infrastructure/resource-quotas';
import { createCertManagerSetup } from './lib/infrastructure/cert-manager-setup';

// Import component functions
import { createMediaFusionDeployment } from './lib/components/mediafusion';
import { createCometDeployment } from './lib/components/comet';
import { createInfrastructureServices } from './lib/components/infrastructure-services';
import { createTraefikHelmChartConfig } from './lib/components/traefik-config';
import { createTraefikIngressRoutes } from './lib/components/traefik-ingressroutes';
import { createServicesIngressRoutes } from './lib/components/services-ingressroutes';
import { createIndexerServices } from './lib/components/indexer-services';
import { createUtilityServices } from './lib/components/utility-services';
import { createAdditionalServices } from './lib/components/additional-services';
import { createVpnGateway } from './lib/components/vpn-gateway';
import { createDevelopmentServices } from './lib/components/development-services';
import { createAiServices } from './lib/components/ai-services';

// Infrastructure Chart - Core cluster resources
export class InfrastructureChart extends Chart {
  constructor(scope: Construct, id: string, props: ChartProps = {}) {
    super(scope, id, props);

    // Create namespaces
    createNamespaces(this);

    // Create resource quotas and limit ranges
    createResourceQuotas(this);
  }
}

// Certificate Manager Chart - TLS certificate management
export class CertManagerChart extends Chart {
  constructor(scope: Construct, id: string, props: ChartProps = {}) {
    super(scope, id, props);

    // Create cert-manager setup job
    createCertManagerSetup(this);
  }
}

// VPN Gateway Chart - Network gateway and VPN services
export class VpnGatewayChart extends Chart {
  constructor(scope: Construct, id: string, props: ChartProps = {}) {
    super(scope, id, props);

    // Create VPN gateway and pod-gateway configuration
    createVpnGateway(this);
  }
}

// Infrastructure Services Chart - Supporting services
export class InfrastructureServicesChart extends Chart {
  constructor(scope: Construct, id: string, props: ChartProps = {}) {
    super(scope, id, props);

    // Create infrastructure services (MongoDB, Redis, etc.)
    createInfrastructureServices(this);
  }
}

// Media Services Chart - Media-related applications
export class MediaServicesChart extends Chart {
  constructor(scope: Construct, id: string, props: ChartProps = {}) {
    super(scope, id, props);

    // Create MediaFusion deployment
    createMediaFusionDeployment(this);

    // Create Comet deployment
    createCometDeployment(this);

    // Create indexer services (Prowlarr, Jackett)
    createIndexerServices(this);

    // Create utility services (Homepage, etc.)
    createUtilityServices(this);

    // Create additional services (SearXNG, Dozzle, Speedtest, etc.)
    createAdditionalServices(this);
  }
}

// Development and Management Chart - Development tools and management interfaces
export class DevelopmentManagementChart extends Chart {
  constructor(scope: Construct, id: string, props: ChartProps = {}) {
    super(scope, id, props);

    // Create development services (Code Server, etc.)
    createDevelopmentServices(this);
  }
}

// AI and Chat Services Chart - AI-powered applications and chat services
export class AiChatServicesChart extends Chart {
  constructor(scope: Construct, id: string, props: ChartProps = {}) {
    super(scope, id, props);

    // Create AI services (LobeChat, etc.)
    createAiServices(this);
  }
}

// Ingress and Routing Chart - Traffic management
export class IngressRoutingChart extends Chart {
  constructor(scope: Construct, id: string, props: ChartProps = {}) {
    super(scope, id, props);

    // Create Traefik HelmChartConfig
    createTraefikHelmChartConfig(this);

    // Create Traefik dashboard IngressRoutes
    createTraefikIngressRoutes(this);

    // Create services IngressRoutes
    createServicesIngressRoutes(this);
  }
}

// Main application
const app = new App();

// Create charts with proper dependencies
const infrastructureChart = new InfrastructureChart(app, 'infrastructure');
const certManagerChart = new CertManagerChart(app, 'cert-manager');
const vpnGatewayChart = new VpnGatewayChart(app, 'vpn-gateway');
// TODO: Add other charts as MyPrecious provides comprehensive services
// const infrastructureServicesChart = new InfrastructureServicesChart(app, 'infrastructure-services');
// const mediaServicesChart = new MediaServicesChart(app, 'media-services');
// const developmentManagementChart = new DevelopmentManagementChart(app, 'development-management');
// const aiChatServicesChart = new AiChatServicesChart(app, 'ai-chat-services');
// const ingressRoutingChart = new IngressRoutingChart(app, 'ingress-routing');

// Import the MyPreciousChart (comprehensive ElfHosted replacement)
import { MyPreciousChart } from './lib/elfhosted/myprecious-chart';

// Create MyPreciousChart (replaces the previous ElfHostedChart)
const myPreciousChart = new MyPreciousChart(app, 'myprecious');

// Set up dependencies following cdk8s best practices
certManagerChart.addDependency(infrastructureChart);
vpnGatewayChart.addDependency(certManagerChart);
//infrastructureServicesChart.addDependency(vpnGatewayChart);
//mediaServicesChart.addDependency(infrastructureServicesChart);
//developmentManagementChart.addDependency(infrastructureServicesChart);
//aiChatServicesChart.addDependency(infrastructureServicesChart);
//ingressRoutingChart.addDependency(mediaServicesChart);
//ingressRoutingChart.addDependency(developmentManagementChart);
//ingressRoutingChart.addDependency(aiChatServicesChart);

// Add MyPrecious chart dependencies (comprehensive replacement for other services)
myPreciousChart.addDependency(vpnGatewayChart);

app.synth();
