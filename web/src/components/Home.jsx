import React from 'react';

export default React.createClass({
  getPair: function() {
    return this.props.pair || [];
  },
  render: function() {
    return <div className="Home">
      <h1>Home</h1>
    </div>;
  }
});
