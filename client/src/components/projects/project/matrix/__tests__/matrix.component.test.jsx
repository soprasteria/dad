import Matrix from '../matrix.component';

import deepFreeze from 'deep-freeze';
import { options, /*getProgressOptions*/ } from '../../../../../modules/services/services.constants';

const defaultServices = {
  name: 'Pipeline d\'intégration continue',
  package: '2. Build',
  services: ['jenkins', 'gitlabci', 'tfs']
};

deepFreeze(defaultServices); // To make Object recursively immuable

describe('Testing getServiceStatus function for a functional service that has 3 different ' +
    'technical services associated',
() => {
  const matrix = new Matrix();
  const oneValidIndicator = {
    indicator1: {
      docktorGroup: 'ProjectA',
      service: 'jenkins',
      status: 'Undetermined'
    }
  };
  deepFreeze(oneValidIndicator);
  describe('With a unique functional service and an existing status', () => {
    const statusToDisplay = matrix.getServiceStatus(defaultServices, oneValidIndicator);
    it('Should return the only existing indicator status', () => {
      expect(statusToDisplay.text).toEqual(oneValidIndicator.indicator1.status);
    });
  });

  describe('With a unique functional service and an unexisting status', () => {
    const oneInvalidIndicator = {
      indicator1: {
        docktorGroup: 'ProjectA',
        service: 'jenkins',
        status: 'Unexisting Status'
      }
    };
    deepFreeze(oneInvalidIndicator);
    const undefinedStatus = matrix.getServiceStatus(defaultServices, oneInvalidIndicator);
    it('Should return an undefined value', () => {
      expect(undefinedStatus).toBeUndefined();
    });
  });

  const twoValidIndicators = {
    ...oneValidIndicator,
    indicator2: {
      docktorGroup: 'ProjectA',
      service: 'gitlabci',
      status: 'Inactive'
    }
  };
  deepFreeze(twoValidIndicators);
  describe('With two functional services and an existing status better than the previous one', () => {
    const statusToDisplay = matrix.getServiceStatus(defaultServices, twoValidIndicators);
    it('Should return the best valid indicator status', () => {
      expect(statusToDisplay.text).toEqual(twoValidIndicators.indicator2.status);
    });
  });

  describe('With another indicator, with a better status but not included in defaultServices', () => {
    const threeNotAllValidIndicators = {
      ...twoValidIndicators,
      indicator3: {
        docktorGroup: 'ProjectA',
        service: 'Unexisting Service',
        status: 'Active'
      }
    };
    deepFreeze(threeNotAllValidIndicators);
    const statusToDisplay = matrix.getServiceStatus(defaultServices, threeNotAllValidIndicators);
    it('Should return the best valid indicator status', () => {
      expect(statusToDisplay.text).toEqual(threeNotAllValidIndicators.indicator2.status);
    });
  });

  describe('With another indicator, with a worse valid status', () => {
    const threeValidIndicators = {
      ...twoValidIndicators,
      indicator3: {
        docktorGroup: 'ProjectA',
        service: 'tfs',
        status: 'Empty'
      }
    };
    deepFreeze(threeValidIndicators);
    const statusToDisplay = matrix.getServiceStatus(defaultServices, threeValidIndicators);
    it('Should return the best valid indicator status', () => {
      expect(statusToDisplay.text).toEqual(threeValidIndicators.indicator2.status);
    });
  });

  describe('With an empty service list and an indicator', () => {
    const emptyServices = {
      name: 'Pipeline d\'intégration continue',
      package: '2. Build',
      services: []
    };
    deepFreeze(emptyServices);
    const statusToDisplay = matrix.getServiceStatus(emptyServices, oneValidIndicator);
    it('Should return an undefined status and not an error', () => {
      expect(statusToDisplay).toBeUndefined();
    });
  });

  describe('With an undefined service list and an indicator', () => {
    const undefinedServices = {
      name: 'Pipeline d\'intégration continue',
      package: '2. Build'
    };
    deepFreeze(undefinedServices);
    const statusToDisplay = matrix.getServiceStatus(undefinedServices, oneValidIndicator);
    it('Should return an undefined status and not an error', () => {
      expect(statusToDisplay).toBeUndefined();
    });
  });
});

const defaultOptions = options;

deepFreeze(defaultOptions);
