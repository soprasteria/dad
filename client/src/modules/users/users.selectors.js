import { containsWithoutAccents } from '../utils/strings';
import { sortby } from '../utils/arrays';

export const sortUsers = (u1, u2) => {
  let comp = 0;
  if (u1.role === 'admin' && (u2.role === 'ri' || u2.role === 'pm')) {
    comp = -1;
  } else if (u1.role === 'ri' && u2.role === 'pm') {
    comp = -1;
  } else if (u1.role === 'ri' && u2.role === 'admin') {
    comp = 1;
  } else if (u1.role === 'pm' && (u2.role === 'admin' || u2.role === 'ri')) {
    comp = 1;
  }
  if (comp === 0) {
    return `${u1.lastName} ${u1.firstName}`.localeCompare(`${u2.lastName} ${u2.firstName}`);
  }
  return comp;
};

export const getFilteredUsers = (users, filterValue) => {
  if (!filterValue || filterValue === '') {
    return Object.values(users).sort(sortUsers);
  } else {
    return Object.values(users).filter((user) => {
      return containsWithoutAccents(JSON.stringify(Object.values(user)), filterValue);
    }).sort(sortUsers);
  }
};

export const getUsersAsOptions = (users) => {
  return [{ value: '', text: 'None' }].concat(users.map((user) => {
    return { value: user.id, text: `${user.lastName.toUpperCase()} ${user.firstName}` } ;
  }).sort(sortby('text')));
};
