import Matrix from '../matrix.component';

const defaultServices = {
  name: 'Pipeline d\'intégration continue',
  package: '2. Build',
  services: [
    'jenkins',
    'gitlabci',
    'tfs'
  ]
};

describe('Testing getServiceStatus function for a functionnal service that has 3 different technical services associated', () => {
  const matrix = new Matrix();
  const OneValidIndicator = {
    indicator1: {
      docktorGroup: 'ProjectA',
      service: 'jenkins',
      status: 'Undetermined'
    }
  };
  describe('With an unique functionnal service (and an existing status)', () => {
    const StatusToDisplay = matrix.getServiceStatus(defaultServices, OneValidIndicator);
    it('Should return the only existing indicator status', () => {
      expect(StatusToDisplay.text).toEqual(OneValidIndicator.indicator1.status);
    });
  });

  describe('With an unique functionnal service (and an unexisting status)', () => {
    const OneInvalidIndicator = {
      indicator1: {
        docktorGroup: 'ProjectA',
        service: 'jenkins',
        status: 'Unexisting Status'
      }
    };
    const UndefinedStatus = matrix.getServiceStatus(defaultServices, OneInvalidIndicator);
    it('Whould return an undefined value', () => {
      expect(UndefinedStatus).toBeUndefined();
    });
  });

  const TwoValidsIndicators = {
    ...OneValidIndicator,
    indicator2: {
      docktorGroup: 'ProjectA',
      service: 'gitlabci',
      status: 'Inactive',
    }
  };
  describe('With Two functionnal services (and an existing status better than the previous one)', () => {
    const StatusToDisplay = matrix.getServiceStatus(defaultServices, TwoValidsIndicators);
    it('Should return the best valid indicator status', () => {
      expect(StatusToDisplay.text).toEqual(TwoValidsIndicators.indicator2.status);
    });
  });

  describe('With another indicator, with a better status but not inclued in defaultServices', () => {
    const ThreeNotAllValidsIndicators = {
      ...TwoValidsIndicators,
      indicator3: {
        docktorGroup: 'ProjectA',
        service: 'Unexisting Service',
        status: 'Active',
      }
    };
    const StatusToDisplay = matrix.getServiceStatus(defaultServices, ThreeNotAllValidsIndicators);
    it('Should return the best valid indicator status', () => {
      expect(StatusToDisplay.text).toEqual(ThreeNotAllValidsIndicators.indicator2.status);
    });
  });

  describe('With another indicator, with a worse valid status', () => {
    const ThreeValidsIndicators = {
      ...TwoValidsIndicators,
      indicator3: {
        docktorGroup: 'ProjectA',
        service: 'tfs',
        status: 'Empty',
      }
    };
    const StatusToDisplay = matrix.getServiceStatus(defaultServices, ThreeValidsIndicators);
    it('Should return the best valid indicator status', () => {
      expect(StatusToDisplay.text).toEqual(ThreeValidsIndicators.indicator2.status);
    });
  });
});

