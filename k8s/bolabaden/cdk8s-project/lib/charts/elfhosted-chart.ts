import { Chart, ChartProps } from 'cdk8s';
import { Construct } from 'constructs';

import { createNamespace } from '../elfhosted/namespace';
import { createStorageResources } from '../elfhosted/storage';
import { createConfigMaps } from '../elfhosted/configmaps';
import { createFilebrowserResources } from '../elfhosted/filebrowser';
import { createGatusResources } from '../elfhosted/gatus';
import { createHomerResources } from '../elfhosted/homer';
import { createPlexResources } from '../elfhosted/plex';
import { createRivenResources } from '../elfhosted/riven';
import { createRcloneFMResources } from '../elfhosted/rclonefm';
import { createZurgResources } from '../elfhosted/zurg';
import { createWizarrResources } from '../elfhosted/wizarr';
import { createKubernetesDashboardResources } from '../elfhosted/kubernetesdashboard';
import { createRcloneUIResources } from '../elfhosted/rcloneui';
import { createTraefikForwardAuthResources } from '../elfhosted/traefikforwardauth';
import { createElfTermResources } from '../elfhosted/elfterm';
import { createStorageHubResources } from '../elfhosted/storagehub';
import { createAdditionalStorageResources } from '../elfhosted/additional-storage';
import { createAdditionalConfigMaps } from '../elfhosted/additional-configmaps';


/**
 * ElfHosted Chart - Contains all ElfHosted resources
 */
export class ElfHostedChart extends Chart {
  constructor(scope: Construct, id: string, props: ChartProps = {}) {
    super(scope, id, props);

    // Create namespace
    createNamespace(this);

    // Create storage resources
    createStorageResources(this);

    // Create additional storage resources
    createAdditionalStorageResources(this);

    // Create ConfigMaps
    createConfigMaps(this);

    // Create additional ConfigMaps
    createAdditionalConfigMaps(this);

    // Create Filebrowser resources
    createFilebrowserResources(this);

    // Create Gatus resources
    createGatusResources(this);

    // Create Homer resources
    createHomerResources(this);

    // Create Plex resources
    createPlexResources(this);

    // Create Riven resources
    createRivenResources(this);

    // Create RcloneFM resources
    createRcloneFMResources(this);

    // Create Zurg resources
    createZurgResources(this);

    // Create Wizarr resources
    createWizarrResources(this);

    // Create Kubernetes Dashboard resources
    createKubernetesDashboardResources(this);

    // Create RcloneUI resources
    createRcloneUIResources(this);

    // Create Traefik Forward Auth resources
    createTraefikForwardAuthResources(this);

    // Create ElfTerm resources
    createElfTermResources(this);

    // Create StorageHub resources
    createStorageHubResources(this);
  }
} 