

export function callAPIMiddleware({dispatch, getstate}) {
  console.log("calling api")
  return next => action => {

    const {
      types,
      method,
      callAPI,
      payload
    } = action;

    console.log(action);

    if (!types) {
      //Normal action: pass it on
      return next(action)
    }

    if (
      !Array.isArray(types) ||
      types.length !== 3 ||
      !types.every(type => typeof type === 'string')
    ) {
      throw new Error('Expected an array of three string types.');
    }

    if (typeof callAPI !== 'function') {
      throw new Error('Expected fetch to be a function.');
    }

    //call api function

    const [requestType, successType, failureType] = types;

    dispatch(Object.assign({}, payload, {
      type: requestType,
    }));



  }
}

/*****
 * code based off of Jared Palmer react-production-starter
 * https://github.com/jaredpalmer/react-production-starter/blob/master/src/middleware/callAPIMiddleware.js
 */
