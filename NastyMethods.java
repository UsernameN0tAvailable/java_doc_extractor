/*
 * Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
 * or more contributor license agreements. Licensed under the Elastic License
 * 2.0 and the Server Side Public License, v 1; you may not use this file except
 * in compliance with, at your election, the Elastic License 2.0 or the Server
 * Side Public License, v 1.
 */
package org.elasticsearch.gradle.testkit;

import org.junit.Assert;
import org.junit.Test;

public class NastyInnerClasses {

	public int method1() {

		if (1 > 2) {
			while(true) {

			}
		}
		return 12;

	}

	int method2() {

		if (1 > 2) {
			while(true) {

			}
		}
		return 12;

	}

	int method3(String[] params, int numba) {

		if (1 > 2) {
			while(true) {

			}
		}
		return 12;

	}

}
