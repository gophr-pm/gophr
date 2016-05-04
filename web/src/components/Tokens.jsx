import React from 'react';

export default React.createClass({
  getPair: function() {
    return this.props.pair || [];
  },
  render: function() {
    return <div className="Tokens">
      <h1>Tokens</h1>
    </div>;
  }
});
