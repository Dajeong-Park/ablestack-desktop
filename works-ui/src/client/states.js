export default {
  /**
   * The Guacamole connection has not yet been attempted.
   *
   * @type String
   */
  IDLE: "IDLE",

  /**
   * The Guacamole connection is being established.
   *
   * @type String
   */
  CONNECTING: "CONNECTING",

  /**
   * The Guacamole connection has been successfully established, and the
   * client is now waiting for receipt of initial graphical data.
   *
   * @type String
   */
  WAITING: "WAITING",

  /**
   * The Guacamole connection has been successfully established, and
   * initial graphical data has been received.
   *
   * @type String
   */
  CONNECTED: "CONNECTED",

  /**
   * The Guacamole connection has terminated successfully. No errors are
   * indicated.
   *
   * @type String
   */
  DISCONNECTED: "DISCONNECTED",

  /**
   * The Guacamole connection has terminated due to an error reported by
   * the client. The associated error code is stored in statusCode.
   *
   * @type String
   */
  CLIENT_ERROR: "CLIENT_ERROR",

  /**
   * The Guacamole connection has terminated due to an error reported by
   * the tunnel. The associated error code is stored in statusCode.
   *
   * @type String
   */
  TUNNEL_ERROR: "TUNNEL_ERROR",

  UNSTABLE: "UNSTABLE",

  //file upload status
  STREAM_STATE: {
    IDLE: "IDLE",
    OPEN: "OPEN",
    CLOSED: "CLOSED",
    ERROR: "ERROR",
  },

  //file upload error status
  ERROR_TYPE: {
    BAD_REQUEST: "BAD_REQUEST",
    INVALID_CREDENTIALS: "INVALID_CREDENTIALS",
    INSUFFICIENT_CREDENTIALS: "INSUFFICIENT_CREDENTIALS",
    INTERNAL_ERROR: "INTERNAL_ERROR",
    NOT_FOUND: "NOT_FOUND",
    PERMISSION_DENIED: "PERMISSION_DENIED",
    STREAM_ERROR: "STREAM_ERROR",
  },
};
