import { initialState, generateEntitiesReducer, generateEntitiesConstants, generateEntitiesActions } from '../entities';

import deepFreeze from 'deep-freeze';

describe('constants generator', () => {
  it('should generate the correctly named constants', () => {
    const name = 'tests';
    const expectedConstants = {
      'INVALID_DELETE_TEST': 'INVALID_DELETE_TEST',
      'INVALID_REQUEST_TEST': 'INVALID_REQUEST_TEST',
      'INVALID_REQUEST_TESTS': 'INVALID_REQUEST_TESTS',
      'INVALID_SAVE_TEST': 'INVALID_SAVE_TEST',
      'RECEIVE_TEST': 'RECEIVE_TEST',
      'RECEIVE_TESTS': 'RECEIVE_TESTS',
      'REQUEST_DELETE_TEST': 'REQUEST_DELETE_TEST',
      'REQUEST_SAVE_TEST': 'REQUEST_SAVE_TEST',
      'REQUEST_TEST': 'REQUEST_TEST',
      'REQUEST_TESTS': 'REQUEST_TESTS',
      'TEST_DELETED': 'TEST_DELETED',
      'TEST_SAVED': 'TEST_SAVED',
    };
    const constants = generateEntitiesConstants(name);
    expect(constants).toEqual(expectedConstants);
  });
});

describe('actions generator', () => {
  it('should generate the correctly named actions', () => {
    const name = 'tests';

    const actions = generateEntitiesActions(name);
    expect(actions).toHaveProperty('requestAll');
    expect(actions).toHaveProperty('receiveSome');
    expect(actions).toHaveProperty('invalidRequest');
    expect(actions).toHaveProperty('requestOne');
    expect(actions).toHaveProperty('receiveOne');
    expect(actions).toHaveProperty('invalidRequestEntity');
    expect(actions).toHaveProperty('requestSave');
    expect(actions).toHaveProperty('saved');
    expect(actions).toHaveProperty('invalidSaveEntity');
    expect(actions).toHaveProperty('requestDelete');
    expect(actions).toHaveProperty('deleted');
    expect(actions).toHaveProperty('invalidDeleteEntity');
  });
});

describe('reducer generator', () => {
  describe('with an unknown action type', () => {
    describe('with a custom initial state', () => {
      it('should not change the state', () => {
        const unknownAction = {
          type: 'UNKNOWN'
        };
        const initialState = {};
        const expectedState = {};

        deepFreeze(unknownAction);
        deepFreeze(initialState);

        const newState = generateEntitiesReducer(initialState, unknownAction, 'tests');

        expect(newState).toEqual(expectedState);
      });
    });

    describe('with the default initial state', () => {
      it('should not change the state', () => {
        const unknownAction = {
          type: 'UNKNOWN'
        };

        deepFreeze(unknownAction);
        deepFreeze(initialState);

        const newState = generateEntitiesReducer(undefined, unknownAction, 'tests');

        expect(newState).toEqual(initialState);
      });
    });
  });

  describe('with the REQUEST_TEST action type', () => {
    it('should set isFetching to true and didInvalidate to false', () => {
      const requestAction = {
        type: 'REQUEST_TESTS'
      };
      const state = {
        data: [1, 2, 3],
        isFetching: false,
        didInvalidate: true
      };
      const expectedState = {
        data: [1, 2, 3],
        isFetching: true,
        didInvalidate: false
      };

      deepFreeze(requestAction);
      deepFreeze(state);

      const newState = generateEntitiesReducer(state, requestAction, 'tests');

      expect(newState).toEqual(expectedState);
    });
  });
});
