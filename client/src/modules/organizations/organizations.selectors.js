import { containsWithoutAccents } from '../utils/strings';

export const getFilteredOrganizations = (organizations, filterValue) => {
  if (!filterValue || filterValue === '') {
    return Object.values(organizations);
  } else {
    return Object.values(organizations).filter(organization => {
      return containsWithoutAccents(JSON.stringify(Object.values(organization)), filterValue);
    });
  }
};

export const getOrganizationsAsOptions = (organizations) => {
  return Object.values(organizations).map(organization => {
    return { value: organization.id, text: organization.name } ;
  });
};
