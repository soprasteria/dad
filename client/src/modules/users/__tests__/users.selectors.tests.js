import { getFilteredUsers, sortUsers } from '../users.selectors';

// Given
let admin = 'admin';
let ri = 'ri';
let pm = 'pm';
let deputy = 'deputy';

describe('I wish to verify the sortUsers function', () => {

  // Given
  let greater = -1;
  let equal = 0;
  let lower = 1;

  it('admin should be equal to admin', () => {

    // Given
    let user1 = {
      role: admin
    };
    let user2 = {
      role: admin
    };

    // When
    let comp = sortUsers(user1, user2);

    // Then
    expect(comp).toBe(equal);
  });

  it('admin should be greater than ri', () => {

    // Given
    let user1 = {
      role: admin
    };
    let user2 = {
      role: ri
    };

    // When
    let comp = sortUsers(user1, user2);

    // Then
    expect(comp).toBe(greater);
  });

  it('admin should be greater than pm', () => {

    // Given
    let user1 = {
      role: admin
    };
    let user2 = {
      role: pm
    };

    // When
    let comp = sortUsers(user1, user2);

    // Then
    expect(comp).toBe(greater);
  });

  it('admin should be greater than deputy', () => {

    // Given
    let user1 = {
      role: admin
    };
    let user2 = {
      role: deputy
    };

    // When
    let comp = sortUsers(user1, user2);

    // Then
    expect(comp).toBe(greater);
  });

  it('ri should be lower than admin', () => {

    // Given
    let user1 = {
      role: ri
    };
    let user2 = {
      role: admin
    };

    // When
    let comp = sortUsers(user1, user2);

    // Then
    expect(comp).toBe(lower);
  });

  it('ri should be equal to ri', () => {

    // Given
    let user1 = {
      role: ri
    };
    let user2 = {
      role: ri
    };

    // When
    let comp = sortUsers(user1, user2);

    // Then
    expect(comp).toBe(equal);
  });

  it('ri should be greater than pm', () => {

    // Given
    let user1 = {
      role: ri
    };
    let user2 = {
      role: pm
    };

    // When
    let comp = sortUsers(user1, user2);

    // Then
    expect(comp).toBe(greater);
  });

  it('ri should be greater than deputy', () => {

    // Given
    let user1 = {
      role: ri
    };
    let user2 = {
      role: deputy
    };

    // When
    let comp = sortUsers(user1, user2);

    // Then
    expect(comp).toBe(greater);
  });

  it('pm should be lower than admin', () => {

    // Given
    let user1 = {
      role: pm
    };
    let user2 = {
      role: admin
    };

    // When
    let comp = sortUsers(user1, user2);

    // Then
    expect(comp).toBe(lower);
  });

  it('pm should be lower than ri', () => {

    // Given
    let user1 = {
      role: pm
    };
    let user2 = {
      role: ri
    };

    // When
    let comp = sortUsers(user1, user2);

    // Then
    expect(comp).toBe(lower);
  });

  it('pm should be equal to pm', () => {

    // Given
    let user1 = {
      role: pm
    };
    let user2 = {
      role: pm
    };

    // When
    let comp = sortUsers(user1, user2);

    // Then
    expect(comp).toBe(equal);
  });

  it('pm should be greater than deputy', () => {

    // Given
    let user1 = {
      role: pm
    };
    let user2 = {
      role: deputy
    };

    // When
    let comp = sortUsers(user1, user2);

    // Then
    expect(comp).toBe(greater);
  });

  it('deputy should be lower than admin', () => {

    // Given
    let user1 = {
      role: deputy
    };
    let user2 = {
      role: admin
    };

    // When
    let comp = sortUsers(user1, user2);

    // Then
    expect(comp).toBe(lower);
  });

  it('deputy should be lower than ri', () => {

    // Given
    let user1 = {
      role: deputy
    };
    let user2 = {
      role: ri
    };

    // When
    let comp = sortUsers(user1, user2);

    // Then
    expect(comp).toBe(lower);
  });

  it('deputy should be lower than pm', () => {

    // Given
    let user1 = {
      role: deputy
    };
    let user2 = {
      role: pm
    };

    // When
    let comp = sortUsers(user1, user2);

    // Then
    expect(comp).toBe(lower);
  });

  it('deputy should be equal to deputy', () => {

    // Given
    let user1 = {
      role: deputy
    };
    let user2 = {
      role: deputy
    };

    // When
    let comp = sortUsers(user1, user2);

    // Then
    expect(comp).toBe(equal);
  });

});


describe('I wish to verify the getFilteredUsers function', () => {

  // Given
  let userAdmin = {
    role: admin
  };
  let userRI = {
    role: ri
  };
  let userPM = {
    role: pm
  };
  let userDeputy = {
    role: deputy
  };

  it('users should be sorted as admin, ri, pm, deputy', () => {

    // Given
    let users = [userPM, userRI, userDeputy, userAdmin];
    let sortedUsersExpected = [userAdmin, userRI, userPM, userDeputy];
    // When
    let sortedUsers = getFilteredUsers(users, false);

    // Then
    expect(sortedUsers).toEqual(sortedUsersExpected);
  });

  it('users should be sorted as admin, admin2, ri, ri2, pm, pm2, deputy, deputy2', () => {

    // Given
    let userAdmin2 = userAdmin;
    let userRI2 = userRI;
    let userPM2 = userPM;
    let userDeputy2 = userDeputy;
    let users = [userRI2, userPM, userDeputy2, userAdmin2, userRI, userDeputy, userPM2, userAdmin];
    let sortedUsersExpected = [userAdmin, userAdmin2, userRI, userRI2, userPM, userPM2, userDeputy, userDeputy2];
    // When
    let sortedUsers = getFilteredUsers(users, false);

    // Then
    expect(sortedUsers).toEqual(sortedUsersExpected);
  });
});
