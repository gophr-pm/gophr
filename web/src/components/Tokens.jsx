import React from 'react';

export default React.createClass({
  getPair: function() {
    return this.props.pair || [];
  },
  render: function() {
    return <div className="Tokens">
      <h2>Tokens</h2>
      <table>
        <th>
          <td>
            Tokens
          </td>
          <td>
            Created
          </td>
          <td>
            delete
          </td>
        </th>
      </table>
    </div>;
  }
});
