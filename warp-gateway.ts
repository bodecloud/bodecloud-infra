import { Construct } from "constructs";
import { App, Chart } from "cdk8s";
import { Deployment, IntOrString, PodSecurityContext } from "cdk8s-plus-26";

export class WarpGatewayChart extends Chart {
  constructor(scope: Construct, id: string) {
    super(scope, id);

    const securityContext: PodSecurityContext = {
      sysctls: [
        {
          name: "net.ipv4.ip_forward",
          value: "1",
        },
        {
          name: "net.ipv4.conf.all.forwarding",
          value: "1",
        },
        {
          name: "net.ipv6.conf.all.forwarding",
          value: "1",
        },
      ],
    };

    const deployment = new Deployment(this, "warp-gateway", {
      replicas: 1,
      containers: [
        {
          image: "your-warp-gateway-image:tag",
          securityContext: {
            privileged: true,
            capabilities: {
              add: ["NET_ADMIN", "NET_RAW"],
            },
          },
        },
      ],
      podSecurityContext: securityContext,
    });
  }
}

const app = new App();
new WarpGatewayChart(app, "warp-gateway");
app.synth();
