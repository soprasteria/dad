import { containsWithoutAccents } from '../utils/strings';

export const getFilteredProjects = (projects, filterValue) => {
  if (!filterValue || filterValue === '') {
    return Object.values(projects);
  } else {
    return Object.values(projects).filter(project => {
      return containsWithoutAccents(JSON.stringify(Object.values(project)), filterValue);
    });
  }
};

export const getProjectsAsOptions = (projects) => {
  return Object.values(projects).map(project => {
    return { value: project.id, name: project.name } ;
  });
};
