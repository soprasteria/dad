export const groupByPackage = (services) => {
  const groupBy = {};
  Object.values(services).forEach((service) => {
    groupBy[service.package] = groupBy[service.package] || [];
    groupBy[service.package].push(service);
  });
  return groupBy;
};
