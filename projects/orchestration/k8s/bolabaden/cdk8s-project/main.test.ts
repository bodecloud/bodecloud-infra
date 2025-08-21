import { Testing } from 'cdk8s';
import { InfrastructureChart } from './main';

describe('Placeholder', () => {
  test('Empty', () => {
    const app = Testing.app();
    const chart = new InfrastructureChart(app, 'test-chart');
    const results = Testing.synth(chart);
    expect(results).toMatchSnapshot();
  });
});
