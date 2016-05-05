import React from 'react';

export default React.createClass({
  getPair: function() {
    return this.props.pair || [];
  },
  render: function() {
    return <div className="Package">
      <h1 class="package-name">Neural-Network</h1>
      <p class="package-description">Packages allows user to use all the Nueral networking</p>
    </div>;
  }
});
