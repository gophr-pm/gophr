import React from 'react';
import TextField from 'material-ui/lib/TextField';

export default React.createClass({
  getPair: function() {
    return this.props.pair || [];
  },
  render: function() {
    return <div className="Email-Edit">
      <h2>Edit email address</h2>
      <p>An email will be sent to the new address to confirm. Another email will also be sent to the old yasabere@gmail.com address to verify that this is intentional.</p>
      <p>
        <label>Password</label>
        <br/>
        <TextField
          type="text"
          id="fullname"
          name="fullname"
          hintText=""
          />
      </p>

      <p>
        <label>New email address</label>
        <br/>
        <TextField
          type="text"
          id="homepage"
          name="homepage"
          hintText=""
          />
      </p>
    </div>;
  }
});
