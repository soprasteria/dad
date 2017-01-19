export const sortby = (attr) => (elem1, elem2) => {
  return elem1[attr].localeCompare(elem2[attr]);
};
