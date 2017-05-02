import { containsWithoutAccents } from '../utils/strings';
import { sortby } from '../utils/arrays';
import { calculateProgress } from '../utils/projects';


export const getFilteredProjects = (projects, entities, filterValue) => {
  if (!filterValue || filterValue === '') {
    return Object.values(projects)
      .sort(sortby('name'));
  } else {
    return Object.values(projects)
      .filter((project) => {
        const projectContains = containsWithoutAccents(JSON.stringify([project.name, project.domain]), filterValue);
        const businessUnit = entities[project.businessUnit] && entities[project.businessUnit].name;
        const serviceCenter = entities[project.serviceCenter] && entities[project.serviceCenter].name;
        const businessUnitContains = containsWithoutAccents(JSON.stringify(businessUnit || ''), filterValue);
        const serviceCenterContains = containsWithoutAccents(JSON.stringify(serviceCenter || ''), filterValue);
        let matchingProjects = false;
        const filteredMatrixGoal = project.matrix.filter((m) => m.goal >= 0);
        const filteredMatrixProgress = project.matrix.filter((m) => m.progress >= 0);
        if (filterValue === 'started') {
          // We keep all projects whose progression is >= 0% AND the projects with no Goals but with a progress status.
          matchingProjects = Math.floor(calculateProgress(project)) >= 0 || filteredMatrixProgress.length > 0 ;
        }
        if (filterValue === 'no goal') {
          // We keep all projects with a progress status but no goal specified.
          matchingProjects = filteredMatrixGoal.length === 0 && filteredMatrixProgress.length > 0 ;
        }
        if (filterValue === 'not started') {
          // We keep all projects that has not been started <=> no progress status and no goal specified
          matchingProjects = filteredMatrixGoal.length === 0 && filteredMatrixProgress.length === 0 ;
        }
        if (filterValue.endsWith('%')) {
          const parsedValue = Number.parseInt(filterValue.substring(0, filterValue.length - 1));
          if (parsedValue || parsedValue === 0) {
            matchingProjects = Math.floor(calculateProgress(project)) >= parsedValue;
          }
        }
        return projectContains || businessUnitContains || serviceCenterContains || matchingProjects;
      })
      .sort(sortby('name'));
  }
};

export const getProjectsAsOptions = (projects) => {
  return Object.values(projects).map((project) => {
    return { value: project.id, name: project.name } ;
  }).sort(sortby('name'));
};
