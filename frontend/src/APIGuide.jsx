import { useEffect } from "react";
import PropTypes from "prop-types";

function APIGuide(){
    return <>
        how to use api
        </>
}

// APIGuide.propTypes = {
//     setIsLoggedIn: PropTypes.func.isRequired
// }

APIGuide.propTypes = {
    setIsLoggedIn: PropTypes.func,
};

export default APIGuide;