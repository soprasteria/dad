import { containsWithoutAccents } from '../utils/strings';

export const getFilteredUsers = (users, filterValue) => {
  if (!filterValue || filterValue === '') {
    return Object.values(users);
  } else {
    return Object.values(users).filter(user => {
      return containsWithoutAccents(JSON.stringify(Object.values(user)), filterValue);
    });
  }
};

export const getUsersAsOptions = (users) => {
  return Object.values(users).map(user => {
    return { value: user.id, name: user.displayName } ;
  });
};
