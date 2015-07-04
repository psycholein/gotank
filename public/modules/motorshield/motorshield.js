// Generated by CoffeeScript 1.9.3
(function() {
  var Motorshield,
    bind = function(fn, me){ return function(){ return fn.apply(me, arguments); }; },
    extend = function(child, parent) { for (var key in parent) { if (hasProp.call(parent, key)) child[key] = parent[key]; } function ctor() { this.constructor = child; } ctor.prototype = parent.prototype; child.prototype = new ctor(); child.__super__ = parent.prototype; return child; },
    hasProp = {}.hasOwnProperty,
    indexOf = [].indexOf || function(item) { for (var i = 0, l = this.length; i < l; i++) { if (i in this && this[i] === item) return i; } return -1; };

  window.Motorshield = Motorshield = (function(superClass) {
    extend(Motorshield, superClass);

    function Motorshield() {
      this.keyUp = bind(this.keyUp, this);
      this.keyDown = bind(this.keyDown, this);
      this.handleKeys = bind(this.handleKeys, this);
      this.initKeyboardControl = bind(this.initKeyboardControl, this);
      this.send = bind(this.send, this);
      this.control = bind(this.control, this);
      this.afterInit = bind(this.afterInit, this);
      return Motorshield.__super__.constructor.apply(this, arguments);
    }

    Motorshield.prototype.afterInit = function() {
      $(this.selector).find('.control').on('click', this.control);
      return this.initKeyboardControl();
    };

    Motorshield.prototype.control = function(e) {
      var control;
      control = $(e.currentTarget).data('event');
      return this.send(control);
    };

    Motorshield.prototype.send = function(control) {
      return this.event.send(this.module, this.name, 'control', {
        value: control
      });
    };

    Motorshield.prototype.initKeyboardControl = function() {
      this.lastControl = null;
      this.keys = [];
      $('body').on('keydown', this.keyDown);
      return $('body').on('keyup', this.keyUp);
    };

    Motorshield.prototype.handleKeys = function(key, pressed) {
      var $control, control, index;
      switch (key) {
        case 37:
          key = 65;
          break;
        case 38:
          key = 87;
          break;
        case 39:
          key = 68;
          break;
        case 40:
          key = 83;
      }
      $control = $(this.selector).find("[data-key='" + (String.fromCharCode(key)) + "']");
      if ($control.length === 0) {
        return;
      }
      control = $control.data('event');
      if (pressed) {
        if (indexOf.call(this.keys, key) < 0) {
          this.keys.push(key);
        }
      } else {
        index = this.keys.indexOf(key);
        if (index > -1) {
          this.keys.splice(index, 1);
        }
        control = "stop";
        if (this.keys.length > 0) {
          key = String.fromCharCode(this.keys.slice(0));
          $control = $(this.selector).find("[data-key='" + key + "']");
          if ($control.length) {
            control = $control.data('event');
          }
        }
      }
      if (this.lastControl !== control) {
        this.send(control);
      }
      return this.lastControl = control;
    };

    Motorshield.prototype.keyDown = function(e) {
      return this.handleKeys(e.which, true);
    };

    Motorshield.prototype.keyUp = function(e) {
      return this.handleKeys(e.which, false);
    };

    return Motorshield;

  })(window.BasicModule);

}).call(this);
