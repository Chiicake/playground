fn f2c(temp: i32) -> i32 {
    (temp - 32) * 5 / 9
}

fn c2f(temp: i32) -> i32 {
    temp * 9 / 5 + 32
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_f2c_freezing_point() {
        // 测试华氏32度转换为摄氏0度
        assert_eq!(f2c(32), 0);
    }

    #[test]
    fn test_f2c_boiling_point() {
        // 测试华氏212度转换为摄氏100度
        assert_eq!(f2c(212), 100);
    }

    #[test]
    fn test_f2c_negative() {
        // 测试负华氏温度转换
        assert_eq!(f2c(-40), -40); // -40华氏度等于-40摄氏度
    }

    #[test]
    fn test_f2c_rounding() {
        // 测试整数除法舍入
        assert_eq!(f2c(50), 10);  // 50°F = 10°C
        assert_eq!(f2c(51), 10);  // 51°F ≈ 10.555°C 向下取整
        assert_eq!(f2c(52), 11);  // 52°F ≈ 11.111°C 向下取整
    }

    #[test]
    fn test_c2f_freezing_point() {
        // 测试摄氏0度转换为华氏32度
        assert_eq!(c2f(0), 32);
    }

    #[test]
    fn test_c2f_boiling_point() {
        // 测试摄氏100度转换为华氏212度
        assert_eq!(c2f(100), 212);
    }

    #[test]
    fn test_c2f_negative() {
        // 测试负摄氏温度转换
        assert_eq!(c2f(-40), -40); // -40摄氏度等于-40华氏度
    }

    #[test]
    fn test_c2f_rounding() {
        // 测试整数除法舍入
        assert_eq!(c2f(10), 50);  // 10°C = 50°F
        assert_eq!(c2f(11), 51);  // 11°C = 51.8°F 向下取整
        assert_eq!(c2f(12), 53);  // 12°C = 53.6°F 向下取整
    }

    #[test]
    fn test_round_trip() {
        // 测试往返转换
        let original_temp = 25;
        let converted = f2c(c2f(original_temp));
        assert_eq!(converted, original_temp);
    }
}