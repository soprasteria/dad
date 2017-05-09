// import constants
import EntitiesConstants from './indicators.constants';
import { generateIndicatorsActions } from '../utils/indicators';


// Change filter
const changeFilter = (filterValue) => {
  return {
    type: IndicatorsConstants.CHANGE_FILTER,
    filterValue
  };
};

export default {
  ...generateIndicatorsActions('indicators'),
  changeFilter
};
