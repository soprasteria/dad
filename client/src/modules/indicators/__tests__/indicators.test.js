import indicatorsActions from '../indicators.actions';
import indicatorsReducer from '../indicators.reducer';

describe('Indicator Reducer', () => {
  const initialState = {
    isFetching: false,
    didInvalidate: true,
    items: {},
    selected: {},
    lastUpdated: undefined
  };
  describe('with the INVALID_INDICATORS action type', () => {
    const items = [];
    const error = {};
    const newState = indicatorsReducer(initialState, indicatorsActions.invalidRequestEntity(items)(error));
    it('Result state should not be the same object as initial state', () => {
      expect(newState).not.toBe(initialState);
    });
  });
  describe('with the REQUEST_INDICATORS action type', () => {
    const newState = indicatorsReducer(initialState, indicatorsActions.requestSome());
    it('Result state should not be the same object as initial state', () => {
      expect(newState).not.toBe(initialState);
    });
    it('should set isFetching to true', () => {
      expect(newState.isFetching).toBe(true);
    });
    it('should set didInvalidate to false', () => {
      expect(newState.didInvalidate).toBe(false);
    });
  });
  describe('with the RECEIVE_INDICATORS action type', () => {
    const items = [{
      id: 'item1',
      docktorGroup: 'ProjectA',
      service: 'jenkins',
      status: 'Undetermined'
    }];
    const expectedState = {
      isFetching: false,
      didInvalidate: false,
      items: {
        item1: {
          id: 'item1',
          docktorGroup: 'ProjectA',
          service: 'jenkins',
          status: 'Undetermined'
        }
      },
      selected: {},
      lastUpdated: undefined
    };
    const newState = indicatorsReducer(initialState, indicatorsActions.receiveSome(items));
    it('Result state should not be the same object as initial state', () => {
      expect(newState).not.toBe(initialState);
    });
    it('should set isFetching to false', () => {
      expect(newState.isFetching).toBe(false);
    });
    it('should set didInvalidate to false', () => {
      expect(newState.didInvalidate).toBe(false);
    });
    it('ResultState should have items equal to ExpectedState items', () => {
      expect(newState.items).toEqual(expectedState.items);
    });
  });
});
