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

export const getByType = (organizations, type) => {
  return organizations.filter(organization => organization.type == type);
};

export const getOrganizationsAsOptions = (organizations) => {
  return organizations.map(organization => {
    let type = organization.type;
    type = type.charAt(0).toUpperCase() + type.substring(1);

    return { value: organization.id, text: `${type}: ${organization.name}` } ;
  });
};
