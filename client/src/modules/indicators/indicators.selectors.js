import { containsWithoutAccents } from '../utils/strings';
import { sortby } from '../utils/arrays';

export const getFilteredIndicators = (indicators, filterValue) => {
  if (!filterValue || filterValue === '') {
    return Object.values(indicators);
  } else {
    return Object.values(indicators).filter((indicator) => {
      return containsWithoutAccents(JSON.stringify(Object.values(indicator)), filterValue);
    });
  }
};
