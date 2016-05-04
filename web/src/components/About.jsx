import React from 'react';

export default React.createClass({
  getPair: function() {
    return this.props.pair || [];
  },
  render: function() {
    return <div className="About">
      <h1>About</h1>
    </div>;
  }
});
