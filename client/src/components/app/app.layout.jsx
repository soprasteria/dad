// React
import React from 'react';

// Components
import NavBar from './navBar.component';


// Style
import 'semantic-ui-css/semantic.min.css';
import './flex.scss';

// App Component
class App extends React.Component {
  render = () => (
    <div className='layout vertical start-justified fill'>
      <NavBar />
      <div className='flex main layout vertical'>
        {this.props.children}
      </div>
    </div>
  );
}
App.propTypes = { children: React.PropTypes.object };

export default App;
