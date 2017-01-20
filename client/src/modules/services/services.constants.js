import { generateEntitiesConstants } from '../utils/entities';

export const options = [
  { value: -1, text: 'N/A', title:'Non applicable' },
  { value: 0, text: '0%', title:'No action launched on the service' },
  { value: 1, text: '20%', title:'Deployed empty by CDK core team' },
  { value: 2, text: '40%', title:'Configured by project team and ready to use' },
  { value: 3, text: '60%', title:'Used by leaders or seniors' },
  { value: 4, text: '80%', title:'Team trained and aware of the benefits' },
  { value: 5, text: '100%', title:'Used fully by the team' },
];

export default {
  ...generateEntitiesConstants('services'),
  CHANGE_FILTER: 'CHANGE_FILTER_SERVICES'
};
