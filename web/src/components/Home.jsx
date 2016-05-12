import React from 'react';

export default React.createClass({
  getPair: function() {
    return this.props.pair || [];
  },
  render: function() {
    return <div className="Home">
      <h2>Popular Gophr Packages</h2>
      <h2>Getting Started</h2>
      <h2>Recently Updated Packages</h2>
    </div>;
  }
});
