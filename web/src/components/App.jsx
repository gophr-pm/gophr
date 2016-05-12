import Color from 'color';
import Radium, { Style } from 'radium'
import React from 'react';
import { Link } from 'react-router';
import AppBar from 'material-ui/lib/app-bar';
import AutoComplete from 'material-ui/lib/auto-complete';

//import styles from '../stylesheets/style.css';

const packages = [
  'EasyAPI',
  'Sockets',
  'NeuralNetwork',
];

export default React.createClass({
  render: function() {

    const style = {
      margin: '0.5em',
      paddingLeft: 0,
      listStyle: 'none',
      backgroundColor: 'yellow'
    };

    return <div className="App">
      <div className="navbar">
        <div className="nav container" >
          <ul>
            <li><Link to="/about">About</Link></li>
            <li><Link to="/support">Support</Link></li>
          </ul>
        </div>
      </div>
      <div className="container">
        <AutoComplete
          floatingLabelText="find GO packages"
          filter={AutoComplete.fuzzyFilter}
          dataSource={packages}
          style={{textAlign:"center"}}
        />
        <div>
          {this.props.children}
        </div>
      </div>
      <div className="footer" style={styles.footer}>
        <div className="container" >
          <ul style={styles.base}>
            <li><Link to="/about">About</Link></li>
            <li><Link to="/support">Support</Link></li>
            <li><Link to="/tokens">Tokens</Link></li>
            <li><Link to="/tutorial">Tutorial</Link></li>
          </ul>
        </div>
      </div>
      <Style
        rules={styles.template}
      />
    </div>;
  }
});


// You can create your style objects dynamically or share them for
// every instance of the component.

var colors = {
  'primary': '#ffd54f',
  'secondary': '#4dd0e1',
  'tertiary': '#2a2730'
};

var styles = {
  app: {
    position: 'relative',
    minHeight: '500px'
  },

  base: {
    color: '#f00',
    height: '40px',
    width: '100%',

    // Adding interactive state couldn't be easier! Add a special key to your
    // style object (:hover, :focus, :active, or @media) with the additional rules.
    ':hover': {
      background: Color(colors['primary']).lighten(0.2).hexString()
    }
  },

  primary: {
    background: '#0074D9'
  },

  warning: {
    background: '#FF4136'
  },

  pullRight: {
    float: 'right',
    marginRight: '-15px'
  },

  footer: {
    backgroundColor: colors['tertiary'],
    color:'white',
    padding: '20px',
    paddingBottom: '70px',
    bottom: '0px',
    width: '100%'
  },

  template: {
    body: {
      margin: 0,
      fontFamily: 'Helvetica Neue, Helvetica, Arial, sans-serif'
    },
    html: {
      background: '#FFF',
      fontSize: '100%'
    },
    mediaQueries: {
      '(min-width: 550px)': {
        html:  {
          fontSize: '120%'
        }
      },
      '(min-width: 1200px)': {
        html:  {
          fontSize: '140%'
        }
      }
    },
    'h1, h2, h3': {
      fontWeight: 'bold'
    },
    '.container': {
      width: '800px',
      paddingRight:'30px',
      paddingLeft:'30px',
      marginRight:'auto',
      marginLeft:'auto',
      mediaQueries: {
        '(min-width: 550px)': {
          html:  {
            width: '525px'
          }
        },
        '(min-width: 960px)': {
          html:  {
            width: '750px'
          }
        },
        '(min-width: 1200px)': {
          html:  {
            width: '1000px'
          }
        }
      },
    },
    '.footer':{
      ':focus': {
        backgroundColor: '#0088FF'
      },
      'a': {
        ':link': {
          color: Color(colors['tertiary']).lighten(0.8).hexString()
        },
        ':visited': {
          color: Color(colors['tertiary']).lighten(0.6).hexString()
        },
        ':hover': {
          color: Color(colors['tertiary']).lighten(0.9).hexString()
        },
        ':active': {
          color: Color(colors['tertiary']).lighten(0.95).hexString()
        }
      }
    }
  },
};
