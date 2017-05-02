import { calculateProgress } from '../projects';

describe('calculate progress of projects', () => {
  describe('with no goals', () => {
    it('should return NaN', () => {
      const project = {
        matrix: []
      };
      expect(calculateProgress(project)).toEqual(NaN);
    });
  }),
  describe('with some positive goals', () => {
    it('should return the mean of the progress for the positive goals', () => {
      const project = {
        matrix: [
          {
            progress: 3,
            goal: 5
          },
          {
            progress: 1,
            goal: 2
          }
        ]
      };
      expect(calculateProgress(project)).toBeCloseTo(55);
    });
  });
});
