import { generateEntitiesConstants } from '../utils/entities';

export const options = [
  { value: -1, text: 'N/A', label: { color: 'black', empty: true, circular: false }, title:'Non applicable' },
  { value: 0, text: '0%', label: { color: 'grey', empty: true, circular: false }, title:'No action launched on the service' },
  { value: 1, text: '20%', label: { color: 'red', empty: true, circular: false }, title:'Deployed empty by CDK core team' },
  { value: 2, text: '40%', label: { color: 'orange', empty: true, circular: false }, title:'Configured by project team and ready to use' },
  { value: 3, text: '60%', label: { color: 'yellow', empty: true, circular: false }, title:'Used by leaders or seniors' },
  { value: 4, text: '80%', label: { color: 'olive', empty: true, circular: false }, title:'Team trained and aware of the benefits' },
  { value: 5, text: '100%', label: { color: 'green', empty: true, circular: false }, title:'Fully used by the team' },
];

export default {
  ...generateEntitiesConstants('services'),
  CHANGE_FILTER: 'CHANGE_FILTER_SERVICES'
};
