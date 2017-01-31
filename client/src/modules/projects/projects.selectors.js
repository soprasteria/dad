import { containsWithoutAccents } from '../utils/strings';
import { sortby } from '../utils/arrays';

export const getFilteredProjects = (projects, filterValue) => {
  if (!filterValue || filterValue === '') {
    return Object.values(projects)
      .sort(sortby('name'));
  } else {
    return Object.values(projects)
      .filter(project => {
        return containsWithoutAccents(JSON.stringify(Object.values(project)), filterValue);
      })
      .sort(sortby('name'));
  }
};

export const getProjectsAsOptions = (projects) => {
  return Object.values(projects).map(project => {
    return { value: project.id, name: project.name } ;
  }).sort(sortby('name'));
};
