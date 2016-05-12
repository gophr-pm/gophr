import React from 'react';

export default React.createClass({
  getPair: function() {
    return this.props.pair || [];
  },
  render: function() {
    return <div className="Profile">
      <h1>Profile</h1>
      <h2>123 Packages by PowItsYaw</h2>
    </div>;
  }
});
