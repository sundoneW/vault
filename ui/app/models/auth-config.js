/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { belongsTo } from '@ember-data/model';

export default class AuthConfigsModel extends Model {
  @belongsTo('auth-method', { async: false, inverse: 'authConfigs', as: 'auth-config' }) backend;
  getHelpUrl(backend) {
    return `/v1/auth/${backend}/config?help=1`;
  }
}
