export const calculateProgress = (project) => {
  const filteredMatrix = project.matrix.filter(m => m.goal >= 0);
  const goals = filteredMatrix.map(m => [m.progress, m.goal])
    .reduce((acc, [progress, goal]) => {
      if (progress === -1) {progress = 0;}
      if (goal === 0) { return acc; }
      const res = acc  + Math.min(progress * 100 / goal, 100);
      return res;
    }, 0);
  return goals / filteredMatrix.length;
};
