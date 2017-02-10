import { containsWithoutAccents } from '../utils/strings';
import { sortby } from '../utils/arrays';

export const getFilteredProjects = (projects, entities, filterValue) => {
  if (!filterValue || filterValue === '') {
    return Object.values(projects)
      .sort(sortby('name'));
  } else {
    return Object.values(projects)
      .filter(project => {
        const projectContains = containsWithoutAccents(JSON.stringify([project.name, project.domain]), filterValue);
        const businessUnit = entities[project.businessUnit] && entities[project.businessUnit].name;
        const serviceCenter = entities[project.serviceCenter] && entities[project.serviceCenter].name;
        const businessUnitContains = containsWithoutAccents(JSON.stringify(businessUnit || ''), filterValue);
        const serviceCenterContains = containsWithoutAccents(JSON.stringify(serviceCenter || ''), filterValue);
        let goalsMatch = false;
        if(filterValue.endsWith('%')) {
          const parsedValue = Number.parseInt(filterValue.substring(0, filterValue.length - 1));
          if (parsedValue) {
            const filteredMatrix = project.matrix.filter(m => m.goal !== -1);
            if (filteredMatrix.length > 0) {
              const goals = filteredMatrix.map(m => [m.progress, m.goal])
                .reduce((acc, [progress, goal]) => {
                  if (progress === -1) {progress = 0;}
                  const res = acc  + Math.min(progress * 100 / goal, 100);
                  return res;
                }, 0);
              goalsMatch = Math.floor(goals / filteredMatrix.length) >= parsedValue;
            }
          }
        }
        return projectContains || businessUnitContains || serviceCenterContains || goalsMatch;
      })
      .sort(sortby('name'));
  }
};

export const getProjectsAsOptions = (projects) => {
  return Object.values(projects).map(project => {
    return { value: project.id, name: project.name } ;
  }).sort(sortby('name'));
};
