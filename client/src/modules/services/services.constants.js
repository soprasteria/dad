import { generateEntitiesConstants } from '../utils/entities';

export const options = [
  {
    value: -1,
    text: 'N/A',
    label: {
      color: 'black',
      empty: true,
      circular: false
    },
    title: 'Not applicable'
  }, {
    value: 0,
    text: '0%',
    label: {
      color: 'grey',
      empty: true,
      circular: false
    },
    title: 'No action launched on the service'
  }, {
    value: 1,
    text: '20%',
    label: {
      color: 'red',
      empty: true,
      circular: false
    },
    title: 'Deployed empty by CDK core team '
  }, {
    value: 2,
    text: '40%',
    label: {
      color: 'orange',
      empty: true,
      circular: false
    },
    title: 'Configured by project team and ready to use'
  }, {
    value: 3,
    text: '60%',
    label: {
      color: 'yellow',
      empty: true,
      circular: false
    },
    title: 'Used by leaders or seniors'
  }, {
    value: 4,
    text: '80%',
    label: {
      color: 'olive',
      empty: true,
      circular: false
    },
    title: 'Team trained and aware of the benefits'
  }, {
    value: 5,
    text: '100%',
    label: {
      color: 'green',
      empty: true,
      circular: false
    },
    title: 'Fully used by the team'
  }
];

export const status = [
  {
    value: 0,
    text: 'Empty',
    color: 'black',
    title: 'The service was never used by the project'
  }, {
    value: 1,
    text: 'Undetermined',
    color: 'red',
    title: 'We cannot determine if the service is active or not'
  }, {
    value: 2,
    text: 'Inactive',
    color: 'orange',
    title: 'The service has not been active recently'
  }, {
    value: 3,
    text: 'Active',
    color: 'green',
    title: 'The service has been recently active'
  }
];

export const priorities = [
  {
    value: 'N/A',
    text: 'N/A',
    title: 'Not applicable'
  }, {
    value: 'P0',
    text: 'P0',
    title: 'High priority'
  }, {
    value: 'P1',
    text: 'P1',
    title: 'Medium priority'
  }, {
    value: 'P2',
    text: 'P2',
    title: 'Low priority'
  }
];

export const deployed = [
  {
    value: 'no',
    text: 'No',
    title: 'The service is not deployed'
  }, {
    value: 'yes',
    text: 'Yes',
    title: 'The service has been deployed'
  }
];

export function getDeployedOptions(deployed, deployedValue, isConnectedUserAdmin) {
  const deployedOptions = [...deployed];

  if (deployedValue === 'Yes' && !isConnectedUserAdmin) {
    deployedOptions[0] = {
      ...deployedOptions[0],
      title: 'Only Admin users can now return back to these values',
      disabled: true
    };
  }

  return deployedOptions;
};

export default {
  ...generateEntitiesConstants('services'),
  CHANGE_FILTER: 'CHANGE_FILTER_SERVICES'
};
