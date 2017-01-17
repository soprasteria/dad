import { createStore, combineReducers, applyMiddleware } from 'redux';
import createLogger from 'redux-logger';
import thunkMiddleware from 'redux-thunk';
import { browserHistory } from 'react-router';
import { routerMiddleware, routerReducer } from 'react-router-redux';

// Reducers
import auth from './modules/auth/auth.reducer';

// Thunks
import AuthThunks from './modules/auth/auth.thunk';

// Configure middlewares
const rMiddleware = routerMiddleware(browserHistory);
let middlewares = [ thunkMiddleware, rMiddleware ];
if (process.env.NODE_ENV !== 'production') {
  // Dev dependencies
  const loggerMiddleware = createLogger();
  middlewares = [ ...middlewares, loggerMiddleware ];
}

// Add the reducer to your store on the `routing` key
const store = createStore(
  combineReducers(
    {
      auth,
      routing: routerReducer,
    }
  ),
  applyMiddleware(...middlewares)
);

const authToken = localStorage.getItem('id_token');
if (authToken) {
  //store.dispatch(AuthThunks.profile());
}

export { store };
